package config

import (
	"context"
<<<<<<<< HEAD:event-service/config/grpc.go
	"goodkarma-event-service/middlewares"
	pb "goodkarma-event-service/proto"
	"goodkarma-event-service/src/repository"
	"goodkarma-event-service/src/service"
========
	"goodkarma-notification-service/middlewares"
	pb "goodkarma-notification-service/proto"
	"goodkarma-notification-service/src/service"
>>>>>>>> 6d7a7a5acaf20b0862bcffd64f136c7266be2e9c:notification-service/config/grpc.go
	"log"
	"net"
	"os"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

func ListenAndServeGrpc() {
	port := os.Getenv("PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	db := InitDatabase()
	eventRepository := repository.NewEventRepository(db)
	notifService := service.NewEventService(eventRepository)

	// Define a custom interceptor for JWT that conditionally skips authentication for register endpoint
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(middlewares.NewInterceptorLogger()),
			grpc_auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
				return middlewares.JWTAuth(ctx)
			}),
		),
	)

	pb.RegisterEventServiceServer(grpcServer, notifService)

	log.Println("\033[36mGRPC server is running on port:", port, "\033[0m")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server gRPC:", err)
	}
}
