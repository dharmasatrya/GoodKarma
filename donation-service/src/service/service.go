// service/borrow_service.go
package service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/dharmasatrya/goodkarma/donation-service/client"
	"github.com/dharmasatrya/goodkarma/donation-service/entity"
	"github.com/dharmasatrya/goodkarma/donation-service/src/repository"
	"github.com/dharmasatrya/goodkarma/payment-service/proto"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"
	eventpb "github.com/dharmasatrya/goodkarma/event-service/proto"
	karmaPb "github.com/dharmasatrya/goodkarma/karma-service/proto"
	paymentpb "github.com/dharmasatrya/goodkarma/payment-service/proto"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type DonationService struct {
	pb.UnimplementedDonationServiceServer
	donationRepository repository.DonationRepository
	paymentClient      *client.PaymentServiceClient
	eventClient        *client.EventServiceClient
	karmaClient        *client.KarmaServiceClient
	messageBroker      MessageBroker
}

// var jwtSecret = []byte("secret")
func NewDonationService(
	donationRepository repository.DonationRepository,
	paymentClient *client.PaymentServiceClient,
	eventClient *client.EventServiceClient,
	karmaClient *client.KarmaServiceClient,
	messageBroker MessageBroker,
) *DonationService {
	return &DonationService{
		donationRepository: donationRepository,
		paymentClient:      paymentClient,
		eventClient:        eventClient,
		karmaClient:        karmaClient,
		messageBroker:      messageBroker,
	}
}

func (s *DonationService) CreateDonation(ctx context.Context, req *pb.CreateDonationRequest) (*pb.CreateDonationResponse, error) {
	// Get metadata and validate auth
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Get claims from context
	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		log.Println("failed to get user claims")
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Println("user_id not found in claims")
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	// Start a MongoDB session for transaction
	session, err := s.donationRepository.StartSession()
	if err != nil {
		log.Println("failed to start session: ", err)
		return nil, status.Errorf(codes.Internal, "failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	// Start transaction
	var donation *entity.Donation
	err = session.StartTransaction()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}

	// Execute operations within transaction
	err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		// 1. Create donation
		newDonation := entity.Donation{
			ID:           primitive.NewObjectID(),
			UserID:       userID,
			EventID:      req.EventId,
			Amount:       req.Amount,
			Status:       "pending", // Set initial status as pending
			DonationType: req.DonationType,
		}

		// 2. If it's a money donation, create invoice first
		if req.DonationType == "uang" {
			token := md.Get("authorization")
			if len(token) == 0 {
				return status.Errorf(codes.Unauthenticated, "authorization token not found")
			}

			// Create outgoing context with token
			outgoingMD := metadata.New(map[string]string{
				"authorization": token[0],
			})
			outgoingCtx := metadata.NewOutgoingContext(sessCtx, outgoingMD)

			// Create invoice using the donation ID we just generated
			_, err := s.paymentClient.Client.CreateInvoice(outgoingCtx, &proto.CreateInvoiceRequest{
				UserId:      userID,
				ExternalId:  newDonation.ID.Hex(),
				Amount:      req.Amount,
				Description: "Goodkarma donation",
			})
			if err != nil {
				log.Println("failed to create invoice: ", err)
				// Roll back transaction by aborting
				if abortErr := session.AbortTransaction(sessCtx); abortErr != nil {
					return status.Errorf(codes.Internal, "failed to abort transaction: %v", abortErr)
				}
				return status.Errorf(codes.Internal, "failed to create invoice: %v", err)
			}
		}

		// 3. Save donation to database
		savedDonation, err := s.donationRepository.CreateDonationWithSession(sessCtx, newDonation)
		if err != nil {
			log.Println("failed to create donation: ", err)
			// Roll back transaction by aborting
			if abortErr := session.AbortTransaction(sessCtx); abortErr != nil {
				return status.Errorf(codes.Internal, "failed to abort transaction: %v", abortErr)
			}
			return status.Errorf(codes.Internal, "error creating donation: %v", err)
		}

		donation = savedDonation

		// 4. Commit the transaction
		if err := session.CommitTransaction(sessCtx); err != nil {
			log.Println("failed to commit transaction: ", err)
			return status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Println("failed to execute transaction: ", err)
		return nil, err
	}

	// Return response
	return &pb.CreateDonationResponse{
		Id:           donation.ID.Hex(),
		UserId:       donation.UserID,
		EventId:      donation.EventID,
		Amount:       donation.Amount,
		Status:       donation.Status,
		DonationType: donation.DonationType,
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

	eventIdToInt, err := strconv.Atoi(updatedDonation.EventID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error converting event id to int")
	}
	eventID := uint32(eventIdToInt)

	event, eventErr := s.eventClient.Client.GetEventById(ctx, &eventpb.Id{Id: eventID})
	if eventErr != nil {
		return nil, status.Errorf(codes.Internal, "error fetching event")
	}

	charged, err1 := s.paymentClient.Client.ChargeFees(ctx, &paymentpb.ChargeFeesRequest{
		UserId: event.UserId,
		Amount: 2000,
	})

	fmt.Println("ngurang", charged)
	if err1 != nil {
		fmt.Println(err)
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

func (s *DonationService) UpdateDonationStatusXendit(ctx context.Context, req *pb.UpdateDonationStatusRequest) (*pb.UpdateDonationStatusResponse, error) {
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

	// get cashback karma point
	if updatedDonation.Status == "COMPLETED" && updatedDonation.DonationType == "uang" {
		_, err := s.karmaClient.Client.CashbackDonation(ctx, &karmaPb.CashbackDonationRequest{
			UserId: updatedDonation.UserID,
			Amount: uint32(updatedDonation.Amount),
		})
		if err != nil {
			log.Printf("(donation-service): error updating karma point from callback xendit: %v", err)
			return nil, status.Errorf(codes.Internal, "error updating karma point")
		}
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

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized to access this user's donations")
	}

	donations, err := s.donationRepository.GetDonationsByUserId(userID)
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
