package routes

import (
	"gateway-service/config"
	"gateway-service/src/controller"
	"gateway-service/src/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	userClient, err := config.InitUserServiceClient()
	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	eventClient, err := config.InitEventServiceClient()
	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	paymentClient, err := config.InitPaymentServiceClient()
	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	donationClient, err := config.InitDonationServiceClient()
	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	// userClient := pb.NewUserServiceClient(userConnection)
	userService := service.NewUserService(userClient)
	userController := controller.NewUserController(userService)

	eventService := service.NewEventService(eventClient)
	eventController := controller.NewEventController(eventService)

	paymentService := service.NewPaymentService(paymentClient)
	paymentController := controller.NewPaymentController(paymentService)

	donationService := service.NewDonationService(donationClient)
	donationController := controller.NewDonationController(donationService)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	user := e.Group("/users")
	{
		user.POST("/register/supporters", userController.RegisterUserSupporter)
		user.POST("/register/coordinators", userController.RegisterUserCoordinator)
		user.POST("/login", userController.Login)
		user.GET("/:id", userController.GetUserById)
		user.GET("/email/verify/:token", userController.VerifyEmail)
	}

	event := e.Group("/events")
	// event.Use(middlewares.RequireAuth)
	// {
	event.POST("", eventController.CreateEvent)
	event.PUT("/:id", eventController.EditEvent)
	event.GET("", eventController.GetAllEvents)
	event.GET("/:id", eventController.GetEventById)
	event.GET("/user", eventController.GetAllEventByUserLogin)
	event.GET("/category/", eventController.GetAllEventByCategory)
	// }

	donation := e.Group("/donations")
	donation.POST("", donationController.CreateDonation)
	donation.PUT("/:id", donationController.UpdateDonationStatus)
	donation.GET("", donationController.GetAllDonationByUser)
	donation.GET("/:event_id", donationController.GetAllDonationByEventId)

	payment := e.Group("/payments")
	payment.GET("/wallets", paymentController.GetWalletByUserId)
	payment.POST("/withdraw", paymentController.Withdraw)
	payment.POST("/invoice", paymentController.XenditInvoiceCallback)
	payment.POST("/xenditcallback/invoice", paymentController.XenditInvoiceCallback)
	payment.POST("/xenditcallback/disbursement", paymentController.XenditDisbursementCallback)

	return e
}
