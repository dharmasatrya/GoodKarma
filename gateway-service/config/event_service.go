package config

import (
	"os"

	pb "github.com/dharmasatrya/goodkarma/event-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEventServiceClient() (pb.EventServiceClient, error) {
	grpcUri := os.Getenv("EVENT_SERVICE_DEV_URI")

	userConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	userClient := pb.NewEventServiceClient(userConnection)

	return userClient, nil
}
