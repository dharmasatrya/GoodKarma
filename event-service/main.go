package main

import (
	"github.com/dharmasatrya/goodkarma/event-service/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.ListenAndServeGrpc()
}
