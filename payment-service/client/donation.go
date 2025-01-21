package client

import (
	"log"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DonationServiceClient struct {
	Client pb.DonationServiceClient
}

func NewDonationServiceClient(donationServiceUrl string) (*DonationServiceClient, error) {
	grpcUri := "localhost:50052"
	donationConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to donation service: %v", err)
		return nil, err
	}

	client := pb.NewDonationServiceClient(donationConnection)
	return &DonationServiceClient{
		Client: client,
	}, nil
}
