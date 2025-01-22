package config

import (
	"os"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitDonationServiceClient() (pb.DonationServiceClient, error) {
	grpcUri := os.Getenv("DONATION_SERVICE_URI")

	donationConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	donationClient := pb.NewDonationServiceClient(donationConnection)

	return donationClient, nil
}
