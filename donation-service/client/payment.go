package client

import (
	"log"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto" // Import user service proto
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentServiceClient struct {
	Client pb.PaymentServiceClient
}

func NewPaymentServiceClient(userServiceUrl string) (*PaymentServiceClient, error) {
	grpcUri := "localhost:50053"
	paymentConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
		return nil, err
	}

	client := pb.NewPaymentServiceClient(paymentConnection)
	return &PaymentServiceClient{
		Client: client,
	}, nil
}
