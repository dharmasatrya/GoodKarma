package main

import (
	"gateway-service/routes"
	"os"
)

func main() {

	router := routes.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Logger.Fatal(router.Start(":" + port))

}
