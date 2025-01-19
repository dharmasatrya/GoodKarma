package main

import (
	"context"
	"log"
	"net"

	"github.com/dharmasatrya/goodkarma/donation-service/middleware"

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

	log.Println("Server is running on port 50052...")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
