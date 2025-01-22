package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/dharmasatrya/goodkarma/payment-service/client"
	"github.com/dharmasatrya/goodkarma/payment-service/middleware"
	"github.com/joho/godotenv"

	"github.com/dharmasatrya/goodkarma/payment-service/config"
	"github.com/dharmasatrya/goodkarma/payment-service/src/repository"
	"github.com/dharmasatrya/goodkarma/payment-service/src/service"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading .env file")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor),
	)

	db, err := config.ConnectionDB(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to db")
	}

	userServiceUrl := os.Getenv("USER_SERVICE_URI_DEV")
	userClient, err := client.NewUserServiceClient(userServiceUrl)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	donationServiceUrl := os.Getenv("DONATION_SERVICE_URI_DEV")
	donationClient, err := client.NewDonationServiceClient(donationServiceUrl)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	eventServiceUrl := os.Getenv("EVENT_SERVICE_URI_DEV")
	eventClient, err := client.NewEventServiceClient(eventServiceUrl)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	paymentRepository := repository.NewPaymentRepository(db)

	paymentService := service.NewPaymentService(paymentRepository, userClient, donationClient, eventClient)
	pb.RegisterPaymentServiceServer(grpcServer, paymentService)

	log.Println("Server is running on port 50053...")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
