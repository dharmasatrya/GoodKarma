package client

import (
	"log"
	"os"

	pb "github.com/dharmasatrya/goodkarma/karma-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type KarmaServiceClient struct {
	Client pb.KarmaServiceClient
}

func NewKarmaServiceClient() (*KarmaServiceClient, error) {
	grpcUri := os.Getenv("KARMA_SERVICE_URI")
	paymentConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to payment service: %v", err)
		return nil, err
	}

	client := pb.NewKarmaServiceClient(paymentConnection)
	return &KarmaServiceClient{
		Client: client,
	}, nil
}
