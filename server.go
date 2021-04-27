package main

import (
	"log"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/raismaulana/ticketing-event/app/config"
	"github.com/raismaulana/ticketing-event/app/delivery/controller"
	"github.com/raismaulana/ticketing-event/app/delivery/middleware"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/repository"
	"github.com/raismaulana/ticketing-event/app/usecase"
	"gorm.io/gorm"
)

var (
	e                     *casbin.Enforcer                 = casbin.NewEnforcer("app/config/casbin-model.conf", "app/config/casbin-policy.csv")
	db                    *gorm.DB                         = config.SetupDatabaseConnection()
	rdb                   *redis.Client                    = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	eventRepository       repository.EventRepository       = repository.NewEventRepository(db)
	userRepository        repository.UserRepository        = repository.NewUserRepository(db)
	transactionRepository repository.TransactionRepository = repository.NewTransactionRepository(db)
	authCase              usecase.AuthCase                 = usecase.NewAuthCase(userRepository)
	eventCase             usecase.EventCase                = usecase.NewEventCase(eventRepository)
	jwtCase               usecase.JWTCase                  = usecase.NewJWTCase()
	redisCase             usecase.RedisCase                = usecase.NewRedisCase(rdb)
	reportCase            usecase.ReportCase               = usecase.NewReportCase(userRepository, eventRepository, transactionRepository)
	transactionCase       usecase.TransactionCase          = usecase.NewTransactionCase(transactionRepository, eventRepository)
	userCase              usecase.UserCase                 = usecase.NewUserCase(userRepository)
	authController        controller.AuthController        = controller.NewAuthController(authCase, jwtCase)
	eventController       controller.EventController       = controller.NewEventController(eventCase, redisCase)
	transactionController controller.TransactionController = controller.NewTransactionController(transactionCase, userCase, redisCase)
	userController        controller.UserController        = controller.NewUserController(userCase, redisCase)
	reportController      controller.ReportController      = controller.NewReportController(reportCase)
)

func main() {
	defer config.CloseDatabaseConnection(db)

	r := gin.Default()
	initRoutes(r)

	r.Run()
}

func initRoutes(r *gin.Engine) {
	guestRoutes := r.Group("api/auth")
	{
		guestRoutes.POST("/login", authController.Login)
		guestRoutes.POST("/register", authController.Register)
	}

	routes := r.Group("api")
	routes.Use(middleware.AuthMiddleware(jwtCase, e))
	{
		userRoutes := routes.Group("/user")
		{
			userRoutes.GET("/", userController.Fetch)
			userRoutes.GET("/:id", userController.GetByID)
			userRoutes.PUT("/update", userController.Update)
			userRoutes.DELETE("/delete/:id", userController.Delete)
		}

		eventRoutes := routes.Group("/event")
		{
			eventRoutes.POST("/insert", eventController.Insert)
			eventRoutes.GET("/", middleware.GetCache(redisCase), eventController.Fetch)
			eventRoutes.GET("/:id", eventController.GetByID)
			eventRoutes.PUT("/update", eventController.Update)
			eventRoutes.DELETE("/delete/:id", eventController.Delete)
		}

		transactionRoutes := routes.Group("/transaction")
		{
			transactionRoutes.POST("/insert", transactionController.Insert)
			transactionRoutes.GET("/", transactionController.Fetch)
			transactionRoutes.GET("/:id", transactionController.GetByID)
			transactionRoutes.PUT("/update", transactionController.Update)
			transactionRoutes.DELETE("/delete/:id", transactionController.Delete)
		}

		customRoutes := routes.Group("")
		{
			customRoutes.GET("/event/available", middleware.GetCache(redisCase), eventController.FetchAvailable)
			customRoutes.POST("/transaction/buy-event", transactionController.BuyEvent)
			// customRoutes.GET("/report/transaction", userController.AllEventReport)
			customRoutes.GET("/report/transaction", func(c *gin.Context) {
				rows, _ := db.Raw("SELECT c.id as cid, e.id as eid, t.id as tid, p.id as pid, p.fullname as pfullname FROM `users` c JOIN event e on c.id = e.creator_id JOIN transaction t on e.id = t.event_id JOIN users p on t.participant_id = p.id WHERE t.status_payment = 'passed'").Rows()
				cols, _ := rows.Columns()
				result := make(map[string]string)
				var results []interface{}
				for rows.Next() {
					columns := make([]string, len(cols))
					columnPointers := make([]interface{}, len(cols))
					for i := range columns {
						columnPointers[i] = &columns[i]
						log.Println(columnPointers[i], " ", columns[i])
					}
					rows.Scan(columnPointers...)
					for i, colName := range cols {
						result[colName] = columns[i]
					}
					results = append(results, result)
				}
				log.Println(rows)
				log.Println(result)
				log.Println(results)
				c.JSON(http.StatusOK, results)
			})
			customRoutes.GET("/report/creator", reportController.FetchAllReportEventByCreator)
			customRoutes.PUT("/transaction/upload", transactionController.UploadReceipt)
			customRoutes.PUT("/transaction/verify", func(c *gin.Context) {
				var verify Verify
				c.ShouldBindJSON(&verify)
				log.Println(verify)
				log.Println(verify.TransactionId)
				log.Println(verify.Status)
				db.Exec("UPDATE transaction SET status_payment = ? WHERE id = ?", verify.Status, verify.TransactionId)
				var user entity.User
				var event entity.Event
				db.Raw("Select p.email FROM users p JOIN transaction t on p.id = t.participant_id WHERE t.`id` = ?", verify.TransactionId).Scan(&user)
				db.Raw("Select e.link_webinar, e.id FROM event e JOIN transaction t ON t.event_id = e.id WHERE t.`id` = ?", verify.TransactionId).Scan(&event)
				if verify.Status == "passed" {
					helper.SendMail(user.Email, "Here We Bring Your Webinar's Link", "we received your payment, here is your link:"+event.LinkWebinar)
				} else if verify.Status == "failed" {
					helper.SendMail(user.Email, "Failed Payment", "Sorry, your payment is invalid:")
					db.Exec("Update event SET quantity = quantity+1 WHERE id = ?", event.ID)
				} else {
					c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("?", "?", helper.EmptyObj{}))
					return
				}
				c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", helper.EmptyObj{}))
			})
		}
	}
}

type Verify struct {
	TransactionId uint64 `form:"transaction_id" json:"transaction_id" binding:"required"`
	Status        string `form:"status" json:"status" binding:"required"`
}
