package service

import (
	"context"
	"gateway-service/dto"
	"log"
	"net/http"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"
	"google.golang.org/grpc/metadata"
)

type DonationService interface {
	CreateDonation(token string, input dto.CreateDonationRequest) (int, *dto.Donation)
	UpdateDonationStatus(token string, input dto.UpdateDonationStatusRequest) (int, *dto.Donation)
	GetAllDonationByUser(token string) (int, []dto.Donation)
	GetAllDonationByEventId(token string, eventID string) (int, []dto.Donation)
}

type donationService struct {
	Client pb.DonationServiceClient
}

func NewDonationService(donationClient pb.DonationServiceClient) *donationService {
	return &donationService{donationClient}
}

func (s *donationService) CreateDonation(token string, input dto.CreateDonationRequest) (int, *dto.Donation) {

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := s.Client.CreateDonation(ctx, &pb.CreateDonationRequest{
		EventId:      input.EventID,
		Amount:       input.Amount,
		Status:       "PENDING",
		DonationType: input.DonationType,
	})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.Donation{
		ID:           res.Id,
		UserID:       res.UserId,
		EventID:      res.EventId,
		Amount:       res.Amount,
		Status:       res.Status,
		DonationType: res.DonationType,
	}

	return http.StatusOK, &response
}

func (s *donationService) UpdateDonationStatus(token string, input dto.UpdateDonationStatusRequest) (int, *dto.Donation) {
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := s.Client.UpdateDonationStatus(ctx, &pb.UpdateDonationStatusRequest{
		Id:     input.ID,
		Status: input.Status,
	})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.Donation{
		ID:           res.Id,
		UserID:       res.UserId,
		EventID:      res.EventId,
		Amount:       res.Amount,
		Status:       res.Status,
		DonationType: res.DonationType,
	}

	return http.StatusOK, &response
}

func (s *donationService) GetAllDonationByUser(token string) (int, []dto.Donation) {
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := s.Client.GetDonationsByUserId(ctx, &pb.GetDonationsByUserIdRequest{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	donations := make([]dto.Donation, len(res.Donations))
	for i, donation := range res.Donations {
		donations[i] = dto.Donation{
			ID:           donation.Id,
			UserID:       donation.UserId,
			EventID:      donation.EventId,
			Amount:       donation.Amount,
			Status:       donation.Status,
			DonationType: donation.DonationType,
		}
	}

	return http.StatusOK, donations
}

func (s *donationService) GetAllDonationByEventId(token string, eventID string) (int, []dto.Donation) {
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := s.Client.GetDonationsByEventId(ctx, &pb.GetDonationsByEventIdRequest{
		EventId: eventID,
	})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	donations := make([]dto.Donation, len(res.Donations))
	for i, donation := range res.Donations {
		donations[i] = dto.Donation{
			ID:           donation.Id,
			UserID:       donation.UserId,
			EventID:      donation.EventId,
			Amount:       donation.Amount,
			Status:       donation.Status,
			DonationType: donation.DonationType,
		}
	}

	return http.StatusOK, donations
}
