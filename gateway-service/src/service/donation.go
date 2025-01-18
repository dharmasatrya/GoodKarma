package service

import (
	"context"
	"gateway-service/dto"
	"log"
	"net/http"

	pb "github.com/dharmasatrya/goodkarma/donation-service/proto"
)

type DonationService interface {
	CreateDonation(input dto.CreateDonationRequest) (int, *dto.Donation)
	UpdateDonationStatus(input dto.UpdateDonationStatusRequest) (int, *dto.Donation)
}

type donationService struct {
	Client pb.DonationServiceClient
}

func NewDonationService(donationClient pb.DonationServiceClient) *donationService {
	return &donationService{donationClient}
}

func (s *donationService) CreateDonation(input dto.CreateDonationRequest) (int, *dto.Donation) {
	res, err := s.Client.CreateDonation(context.Background(), &pb.CreateDonationRequest{})
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

func (s *donationService) UpdateDonationStatus(input dto.UpdateDonationStatusRequest) (int, *dto.Donation) {
	res, err := s.Client.UpdateDonationStatus(context.Background(), &pb.UpdateDonationStatusRequest{})
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

//get donation by user
//get donation by event
