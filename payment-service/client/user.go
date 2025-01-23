// payment-service/client/user_client.go
package client

import (
	"crypto/tls"
	"log"
	"os"

	pb "github.com/dharmasatrya/goodkarma/user-service/proto" // Import user service proto
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type UserServiceClient struct {
	Client pb.UserServiceClient
}

func NewUserServiceClient(userServiceUrl string) (*UserServiceClient, error) {
	grpcUri := os.Getenv("USER_SERVICE_URI")
	userConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
		return nil, err
	}

	client := pb.NewUserServiceClient(userConnection)
	return &UserServiceClient{
		Client: client,
	}, nil
}
