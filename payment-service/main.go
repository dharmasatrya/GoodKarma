package main

import (
	"context"
	"log"
	"net"

	"github.com/dharmasatrya/goodkarma/payment-service/middleware"

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

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor),
	)

	// grpcServer := grpc.NewServer()

	db, err := config.ConnectionDB(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to db")
	}

	paymentRepository := repository.NewPaymentRepository(db)

	paymentService := service.NewPaymentService(paymentRepository)
	pb.RegisterPaymentServiceServer(grpcServer, paymentService)

	log.Println("Server is running on port 50053...")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
