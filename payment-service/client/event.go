package client

import (
	"log"
	"os"

	pb "github.com/dharmasatrya/goodkarma/event-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventServiceClient struct {
	Client pb.EventServiceClient
}

func NewEventServiceClient(eventServiceUrl string) (*EventServiceClient, error) {
	// grpcUri := "localhost:50055"
	grpcUri := os.Getenv("EVENT_SERVICE_URI")
	eventConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to event service: %v", err)
		return nil, err
	}

	client := pb.NewEventServiceClient(eventConnection)
	return &EventServiceClient{
		Client: client,
	}, nil
}
