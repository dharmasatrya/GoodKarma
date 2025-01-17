package config

import (
	"fmt"
	"log"
	"os"

	pb "gateway-service/pb/user"

	"google.golang.org/grpc"
)

func InitUserServiceClient() (*grpc.ClientConn, pb.UserServiceClient) {
	// Without TLS (use this if the server does not support TLS)
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_URI"), grpc.WithInsecure())
	if err != nil {
		log.Fatal(fmt.Sprintf("Error connection grcp to user-service: %v", err))
	}
	return conn, pb.NewUserServiceClient(conn)
}
