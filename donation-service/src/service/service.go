// service/borrow_service.go
package service

import (
	"context"

	"github.com/dharmasatrya/goodkarma/donation-service/entity"
	"github.com/dharmasatrya/goodkarma/donation-service/src/repository"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"

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
		DonationType: req.DonationType,
	}

	res, err := s.donationRepository.CreateDonation(donation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	return &pb.CreateDonationResponse{
		Id:           res.ID.Hex(),
		UserId:       res.UserID,
		EventId:      res.EventID,
		Amount:       res.Amount,
		Status:       res.Status,
		DonationType: res.DonationType,
	}, nil
}

func (s *DonationService) UpdateDonationStatus(ctx context.Context, req *pb.UpdateDonationStatusRequest) (*pb.UpdateDonationStatusResponse, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid donation ID format")
	}

	donation := entity.Donation{
		ID:     objectID,
		Status: req.Status,
	}

	updatedDonation, err := s.donationRepository.UpdateDonationStatus(donation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating donation status: %v", err)
	}

	return &pb.UpdateDonationStatusResponse{
		Id:           updatedDonation.ID.Hex(),
		UserId:       updatedDonation.UserID,
		EventId:      updatedDonation.EventID,
		Amount:       updatedDonation.Amount,
		Status:       updatedDonation.Status,
		DonationType: updatedDonation.DonationType,
	}, nil
}

func (s *DonationService) GetDonationsByUserId(ctx context.Context, req *pb.GetDonationsByUserIdRequest) (*pb.GetDonationsByUserIdResponse, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Get claims from context
	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Verify that the requested user_id matches the token's user_id
	userID, ok := claims["user_id"].(string)
	if !ok || userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized to access this user's donations")
	}

	donations, err := s.donationRepository.GetDonationsByUserId(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching donations: %v", err)
	}

	var pbDonations []*pb.Donation
	for _, d := range donations {
		pbDonations = append(pbDonations, &pb.Donation{
			Id:           d.ID.Hex(),
			UserId:       d.UserID,
			EventId:      d.EventID,
			Amount:       d.Amount,
			Status:       d.Status,
			DonationType: d.DonationType,
		})
	}

	return &pb.GetDonationsByUserIdResponse{
		Donations: pbDonations,
	}, nil
}

func (s *DonationService) GetDonationsByEventId(ctx context.Context, req *pb.GetDonationsByEventIdRequest) (*pb.GetDonationsByEventIdResponse, error) {
	// No auth check needed as this is a public endpoint
	donations, err := s.donationRepository.GetDonationsByEventId(req.EventId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching donations: %v", err)
	}

	var pbDonations []*pb.Donation
	for _, d := range donations {
		pbDonations = append(pbDonations, &pb.Donation{
			Id:           d.ID.Hex(),
			UserId:       d.UserID,
			EventId:      d.EventID,
			Amount:       d.Amount,
			Status:       d.Status,
			DonationType: d.DonationType,
		})
	}

	return &pb.GetDonationsByEventIdResponse{
		Donations: pbDonations,
	}, nil
}
