package client

import (
	"crypto/tls"
	"log"
	"os"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DonationServiceClient struct {
	Client pb.DonationServiceClient
}

func NewDonationServiceClient(donationServiceUrl string) (*DonationServiceClient, error) {
	grpcUri := os.Getenv("DONATION_SERVICE_URI")
	donationConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		log.Fatalf("Failed to connect to donation service: %v", err)
		return nil, err
	}

	client := pb.NewDonationServiceClient(donationConnection)
	return &DonationServiceClient{
		Client: client,
	}, nil
}
