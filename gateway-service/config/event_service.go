package config

import (
	"crypto/tls"
	"os"

	pb "github.com/dharmasatrya/goodkarma/event-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func InitEventServiceClient() (pb.EventServiceClient, error) {
	eventConnection, err := grpc.Dial(os.Getenv("EVENT_SERVICE_URI"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))

	if err != nil {
		return nil, err
	}

	eventClient := pb.NewEventServiceClient(eventConnection)

	return eventClient, nil
}
