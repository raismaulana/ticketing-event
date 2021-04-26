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
	transactionCase       usecase.TransactionCase          = usecase.NewTransactionCase(transactionRepository, eventRepository)
	userCase              usecase.UserCase                 = usecase.NewUserCase(userRepository)
	authController        controller.AuthController        = controller.NewAuthController(authCase, jwtCase)
	eventController       controller.EventController       = controller.NewEventController(eventCase, redisCase)
	transactionController controller.TransactionController = controller.NewTransactionController(transactionCase, redisCase)
	userController        controller.UserController        = controller.NewUserController(userCase, redisCase)
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
			customRoutes.GET("/report/creator", func(c *gin.Context) {
				creator_id := c.MustGet("user_id")
				var reportCreator []ReportCreator
				var detailWebinar DetailWebinar

				tx := db.Raw("SELECT e.*, SUM(t.amount) `total_amount`, COUNT(t.participant_id) `total_participant` FROM `transaction` t JOIN event e on e.id = t.event_id WHERE e.creator_id = 12 AND t.status_payment = 'passed' AND e.event_end_date <= NOW() GROUP BY t.id ORDER BY t.id", creator_id).Scan(&detailWebinar)

				log.Println(tx.Rows())
				log.Println(detailWebinar)
				log.Println(reportCreator)

				// rows, _ := db.Raw("SELECT c.id as cid, e.id as eid, e.event_end_date, t.id as tid, p.id as pid, p.fullname as pfullname FROM `users` c JOIN event e on c.id = e.creator_id JOIN transaction t on e.id = t.event_id JOIN users p on t.participant_id = p.id WHERE e.creator_id = ? AND t.status_payment = 'passed' AND e.event_end_date <= NOW()", creator_id).Rows()

				// for sum.Next() {
				// 	reportCreator[i].jumlah_participant = sum
				// }
				c.JSON(http.StatusOK, detailWebinar)
			})
		}
	}
}

type ReportCreator struct {
	DetailWebinar DetailWebinar
	Participant   []entity.User
}

type DetailWebinar struct {
	Event            entity.Event `gorm:"embedded"`
	TotalAmount      float64      `gorm:"->" json:"total_amount"`
	TotalParticipant int          `gorm:"->" json:"total_participant"`
}
