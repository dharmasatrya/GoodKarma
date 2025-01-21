package config

import (
	"github.com/dharmasatrya/goodkarma/event-service/src/repository"

	"log"
	"net"
	"os"

	"github.com/dharmasatrya/goodkarma/event-service/middlewares"
	pb "github.com/dharmasatrya/goodkarma/event-service/proto"
	"github.com/dharmasatrya/goodkarma/event-service/src/service"

	"google.golang.org/grpc"
)

func ListenAndServeGrpc() {
	port := os.Getenv("PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	// Define a custom interceptor for JWT that conditionally skips authentication for register endpoint
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares.UnaryAuthInterceptor),
	)

	db := InitDatabase()
	eventRepository := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepository)

	pb.RegisterEventServiceServer(grpcServer, eventService)

	log.Println("\033[36mGRPC server is running on port:", port, "\033[0m")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server gRPC:", err)
	}
}
