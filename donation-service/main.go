package main

import (
	"context"
	"log"
	"net"

	"github.com/dharmasatrya/goodkarma/donation-service/client"
	"github.com/dharmasatrya/goodkarma/donation-service/middleware"
	"github.com/joho/godotenv"

	"github.com/dharmasatrya/goodkarma/donation-service/config"
	"github.com/dharmasatrya/goodkarma/donation-service/src/repository"
	"github.com/dharmasatrya/goodkarma/donation-service/src/service"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":50052")
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

	paymentServiceUrl := "localhost:50053"
	paymentClient, err := client.NewPaymentServiceClient(paymentServiceUrl)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	eventServiceUrl := "localhost:50053"
	eventClient, err := client.NewEventServiceClient(eventServiceUrl)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	donationRepository := repository.NewDonationRepository(db)

	conn, mbChan := config.InitMessageBroker()
	defer conn.Close()

	messageBrokerService := service.NewMessageBroker(mbChan)

	donationService := service.NewDonationService(donationRepository, paymentClient, eventClient, messageBrokerService)
	pb.RegisterDonationServiceServer(grpcServer, donationService)

	log.Println("Server is running on port 50052...")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
