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
	// Initialize gRPC client connection
	userClientConn, userClient := config.InitUserServiceClient()
	defer userClientConn.Close() // Close the connection when the application exits

	userService := service.NewUserService(userClient)
	userController := controller.NewUserController(userService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	user := e.Group("/users")
	{
		// register := users.Group("/register")
		// {
		// 	register.POST("/penyelenggara", userController.RegisterPenyelenggara)
		// 	register.POST("/donatur", userController.RegisterDonatur)
		// }

		user.POST("/login", userController.LoginUser)
		// user.PUT("", userController.EditProfile)
		// user.DELETE("", userController.DeleteProfile)
	}

	// event := e.Group("/events")
	// event.Use(middlewares.RequireAuth)
	// {
	// 	event.POST("", eventController.CreateEvent)
	// 	event.PUT("/:id", eventController.EditEvent)
	// 	event.GET("", eventController.GetAllEvent)
	// 	event.GET("/:id", eventController.GetEventById)
	// 	event.GET("/", eventController.GetAllEventByUserLogin)
	// 	event.GET("/:category", eventController.GetAllEventByCategory)
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
