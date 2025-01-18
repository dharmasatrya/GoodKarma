package main

import (
	"goodkarma-event-service/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.ListenAndServeGrpc()
}
