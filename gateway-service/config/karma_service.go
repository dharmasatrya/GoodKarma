package config

import (
	"os"

	pb "github.com/dharmasatrya/goodkarma/karma-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitKarmaServiceClient() (pb.KarmaServiceClient, error) {
	grpcUri := os.Getenv("KARMA_SERVICE_URI")

	karmaConnection, err := grpc.NewClient(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	karmaClient := pb.NewKarmaServiceClient(karmaConnection)

	return karmaClient, nil
}
