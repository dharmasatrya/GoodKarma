package main

import (
	"context"
	"log"
	"net"

	"github.com/dharmasatrya/goodkarma/user-service/config"
	pb "github.com/dharmasatrya/goodkarma/user-service/proto"
	"github.com/dharmasatrya/goodkarma/user-service/repository"
	"github.com/dharmasatrya/goodkarma/user-service/service"

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
