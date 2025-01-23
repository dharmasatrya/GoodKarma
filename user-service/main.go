package main

import (
	"context"
	"log"
	"net"
	"os"

	karmaPb "github.com/dharmasatrya/goodkarma/karma-service/proto"
	paymentPb "github.com/dharmasatrya/goodkarma/payment-service/proto"
	pb "github.com/dharmasatrya/goodkarma/user-service/proto"

	"github.com/dharmasatrya/goodkarma/user-service/config"
	"github.com/dharmasatrya/goodkarma/user-service/repository"
	"github.com/dharmasatrya/goodkarma/user-service/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	godotenv.Load()

	listen, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db, err := config.Connect(context.Background())

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	server := grpc.NewServer()

	conn, mbChan := config.InitMessageBroker()
	defer conn.Close()

	messageBrokerService := service.NewMessageBroker(mbChan)

	// Get payment service URI
	paymentServiceURI := os.Getenv("PAYMENT_SERVICE_URI")
	if paymentServiceURI == "" {
		log.Fatalf("PAYMENT_SERVICE_URI is not set")
	}

	// Create gRPC connection
	grpcConn, err := grpc.NewClient(paymentServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}
	defer grpcConn.Close()

	// Initialize payment client
	paymentClient := paymentPb.NewPaymentServiceClient(grpcConn)

	// Initialize karma client
	karmaServiceURI := os.Getenv("KARMA_SERVICE_URI")

	grpcConn, err = grpc.NewClient(karmaServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to karma service: %v", err)
	}
	defer grpcConn.Close()

	karmaClient := karmaPb.NewKarmaServiceClient(grpcConn)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, messageBrokerService, paymentClient, karmaClient)

	pb.RegisterUserServiceServer(server, userService)

	log.Println("Server is running on port: 50051")

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
