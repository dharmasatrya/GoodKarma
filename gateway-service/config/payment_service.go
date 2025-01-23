package config

import (
	"crypto/tls"
	"os"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func InitPaymentServiceClient() (pb.PaymentServiceClient, error) {
	grpcUri := os.Getenv("PAYMENT_SERVICE_URI")

	paymentConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))

	if err != nil {
		return nil, err
	}

	paymentClient := pb.NewPaymentServiceClient(paymentConnection)

	return paymentClient, nil
}
