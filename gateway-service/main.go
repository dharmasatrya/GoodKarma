package main

import (
	"gateway-service/routes"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := routes.NewRouter()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router.Logger.Fatal(router.Start(":" + port))

}
