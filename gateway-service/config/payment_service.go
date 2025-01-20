package config

import (
	"os"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitPaymentServiceClient() (pb.PaymentServiceClient, error) {
	grpcUri := os.Getenv("PAYMENT_SERVICE_DEV_URI")

	paymentConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	paymentClient := pb.NewPaymentServiceClient(paymentConnection)

	return paymentClient, nil
}
