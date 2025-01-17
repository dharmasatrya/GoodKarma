package main

import (
	"context"
	"goodkarma-donation-service/middleware"
	"log"
	"net"

	"goodkarma-donation-service/config"
	"goodkarma-donation-service/src/repository"
	"goodkarma-donation-service/src/service"

	pb "goodkarma-donation-service/proto"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor),
	)

	// grpcServer := grpc.NewServer()

	db, err := config.ConnectionDB(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to db")
	}

	donationRepository := repository.NewDonationRepository(db)

	donationService := service.NewDonationService(donationRepository)
	pb.RegisterDonationServiceServer(grpcServer, donationService)

	log.Println("Server is running on port 50054...")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
