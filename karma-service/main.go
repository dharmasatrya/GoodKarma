package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/dharmasatrya/goodkarma/karma-service/middleware"
	pb "github.com/dharmasatrya/goodkarma/karma-service/proto"

	"github.com/dharmasatrya/goodkarma/karma-service/config"
	"github.com/dharmasatrya/goodkarma/karma-service/repository"
	"github.com/dharmasatrya/goodkarma/karma-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	listen, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db, err := config.Connect(context.Background())

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor),
	)

	karmaRepository := repository.NewKarmaRepository(db)
	karmaService := service.NewKarmaService(karmaRepository)

	pb.RegisterKarmaServiceServer(server, karmaService)

	log.Printf("Server is running on port: %v", port)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
