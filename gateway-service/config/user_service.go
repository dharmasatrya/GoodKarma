package config

import (
	"os"

	pb "github.com/dharmasatrya/goodkarma/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserServiceClient() (pb.UserServiceClient, error) {
	grpcUri := os.Getenv("USER_SERVICE_URI")

	userConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	userClient := pb.NewUserServiceClient(userConnection)

	return userClient, nil
}
