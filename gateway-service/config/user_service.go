package config

import (
	"crypto/tls"
	"log"
	"os"

	pb "github.com/dharmasatrya/goodkarma/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func InitUserServiceClient() (pb.UserServiceClient, error) {
	userConnection, err := grpc.Dial(os.Getenv("USER_SERVICE_URI"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}

	userClient := pb.NewUserServiceClient(userConnection)

	return userClient, nil
}
