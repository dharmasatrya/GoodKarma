package client

import (
	"crypto/tls"
	"log"
	"os"

	pb "github.com/dharmasatrya/goodkarma/event-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type EventServiceClient struct {
	Client pb.EventServiceClient
}

func NewEventServiceClient(eventServiceUrl string) (*EventServiceClient, error) {
	grpcUri := os.Getenv("EVENT_SERVICE_URI")
	eventConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		log.Fatalf("Failed to connect to event service: %v", err)
		return nil, err
	}

	client := pb.NewEventServiceClient(eventConnection)
	return &EventServiceClient{
		Client: client,
	}, nil
}
