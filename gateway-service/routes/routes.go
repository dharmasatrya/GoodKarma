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

	userService := service.NewUserService(userClient)
	userController := controller.NewUserController(userService)

	eventService := service.NewEventService(eventClient)
	eventController := controller.NewEventController(eventService)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	user := e.Group("/users")
	{
		user.POST("/register/supporters", userController.RegisterUserSupporter)
		user.POST("/register/coordinators", userController.RegisterUserCoordinator)
		user.POST("/login", userController.Login)
		user.GET("/:id", userController.GetUserById)
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

	// donation := e.Group("/donations")
	// {
	// 	donation.POST("", donationController.CreateDonation)
	// 	donation.PUT("/:id", donationController.UpdateStatus)
	// 	donation.GET("", donationController.GetAllDonationByUserLogin)
	// 	donation.GET("/:id", donationController.GetAllDonationByEventId)
	// }

	// payment := e.Group("/payments")
	// payment.GET("/wallets/:id", paymentController.GetWallet)
	// payment.POST("/withdraw", paymentController.Withdraw)
	// payment.POST("/invoice", paymentController.CreateInvoice)

	return e
}
