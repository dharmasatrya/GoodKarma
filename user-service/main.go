package main

import (
	"context"
	"goodkarma-user-service/config"
	pb "goodkarma-user-service/proto"
	"goodkarma-user-service/repository"
	"goodkarma-user-service/service"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	listen, err := net.Listen("tcp", ":50052")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db, err := config.Connect(context.Background())

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	server := grpc.NewServer()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	pb.RegisterUserServiceServer(server, userService)

	log.Println("Server is running on port: 50052")

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
