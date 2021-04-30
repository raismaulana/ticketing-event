package main

import (
	"time"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/raismaulana/ticketing-event/app/config"
	"github.com/raismaulana/ticketing-event/app/delivery/background"
	"github.com/raismaulana/ticketing-event/app/delivery/controller"
	"github.com/raismaulana/ticketing-event/app/delivery/middleware"
	"github.com/raismaulana/ticketing-event/app/repository"
	"github.com/raismaulana/ticketing-event/app/usecase"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	db                    *gorm.DB                         = config.SetupDatabaseConnection()
	amqpConn              *amqp.Connection                 = config.SetupRabbitmqConnection()
	e                     *casbin.Enforcer                 = casbin.NewEnforcer("app/config/casbin-model.conf", "app/config/casbin-policy.csv")
	rdb                   *redis.Client                    = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	eventRepository       repository.EventRepository       = repository.NewEventRepository(db)
	userRepository        repository.UserRepository        = repository.NewUserRepository(db)
	transactionRepository repository.TransactionRepository = repository.NewTransactionRepository(db)
	backgroundCase        usecase.BackgroundCase           = usecase.NewBackgroundCase(transactionRepository, eventRepository, userRepository)
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
	backgroundTask        background.BackgroundTask        = background.NewBackgroundTask(backgroundCase, amqpConn)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	defer amqpConn.Close()

	s := gocron.NewScheduler(time.Local)
	initBackgroundTask(s)
	defer s.Stop()

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
			userRoutes.GET("/", middleware.GetCache(redisCase), userController.Fetch)
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
			eventRoutes.GET("/available", middleware.GetCache(redisCase), eventController.FetchAvailable)
		}

		transactionRoutes := routes.Group("/transaction")
		{
			transactionRoutes.POST("/insert", transactionController.Insert)
			transactionRoutes.GET("/", middleware.GetCache(redisCase), transactionController.Fetch)
			transactionRoutes.GET("/:id", transactionController.GetByID)
			transactionRoutes.PUT("/update", transactionController.Update)
			transactionRoutes.DELETE("/delete/:id", transactionController.Delete)
			transactionRoutes.POST("/buy-event", transactionController.BuyEvent)
			transactionRoutes.PUT("/upload", transactionController.UploadReceipt)
			transactionRoutes.PUT("/verify", transactionController.VerifyPayment)
		}

		customRoutes := routes.Group("")
		{
			customRoutes.GET("/report/transaction", middleware.GetCache(redisCase), reportController.FetchAllReportUserBoughtEvent)
			customRoutes.GET("/report/creator", reportController.FetchAllReportEventByCreator)
		}
	}
}

func initBackgroundTask(s *gocron.Scheduler) {
	go backgroundTask.ListenerReminderPayment()
	go backgroundTask.ListenerPromotionEvent()
	s.Every(1).Day().StartAt(time.Date(2021, time.April, 30, 5, 45, 00, 0, time.Local)).Tag("Reminder").Do(backgroundTask.SendReminderPayment)
	s.Every(1).Day().StartAt(time.Date(2021, time.April, 30, 10, 45, 00, 0, time.Local)).Tag("Promotion").Do(backgroundTask.SendPromotionEvent)
	s.StartAsync()
}
