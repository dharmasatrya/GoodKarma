// service/borrow_service.go
package service

import (
	"context"
	"goodkarma-donation-service/entity"
	"goodkarma-donation-service/src/repository"

	pb "goodkarma-donation-service/proto"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type DonationService struct {
	pb.UnimplementedDonationServiceServer
	donationRepository repository.DonationRepository
}

// var jwtSecret = []byte("secret")

func NewDonationService(donationRepository repository.DonationRepository) *DonationService {
	return &DonationService{
		donationRepository: donationRepository,
	}
}

func (s *DonationService) CreateDonation(ctx context.Context, req *pb.CreateDonationRequest) (*pb.CreateDonationResponse, error) {

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Get claims from context that was set in auth middleware
	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	donation := entity.Donation{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		EventID:      req.EventId,
		Amount:       req.Amount,
		Status:       req.Status,
		DonationType: req.DonationTypeId,
	}

	res, err := s.donationRepository.CreateDonation(donation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	return &pb.CreateDonationResponse{
		Id:             res.ID.Hex(),
		UserId:         res.UserID,
		EventId:        res.EventID,
		Amount:         res.Amount,
		Status:         res.Status,
		DonationTypeId: res.DonationType,
	}, nil
}
