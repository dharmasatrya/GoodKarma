package client

import (
	"crypto/tls"
	"log"
	"os"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto" // Import user service proto
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type PaymentServiceClient struct {
	Client pb.PaymentServiceClient
}

func NewPaymentServiceClient(userServiceUrl string) (*PaymentServiceClient, error) {
	grpcUri := os.Getenv("PAYMENT_SERVICE_URI")
	paymentConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		log.Fatalf("Failed to connect to payment service: %v", err)
		return nil, err
	}

	client := pb.NewPaymentServiceClient(paymentConnection)
	return &PaymentServiceClient{
		Client: client,
	}, nil
}
