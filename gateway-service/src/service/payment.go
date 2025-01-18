package service

import (
	"context"
	"gateway-service/dto"
	"log"
	"net/http"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto"
)

type PaymentService interface {
	Withdraw(input dto.WithdrawRequest) (int, *dto.Wallet)
	CreateInvoice(input dto.CreateInvoiceRequest) (int, *dto.CreateInvoiceResponse)
}

type paymentService struct {
	Client pb.PaymentServiceClient
}

func NewPaymentService(paymentClient pb.PaymentServiceClient) *paymentService {
	return &paymentService{paymentClient}
}

func (u *paymentService) Withdraw(input dto.WithdrawRequest) (int, *dto.Wallet) {
	res, err := u.Client.Withdraw(context.Background(), &pb.WithdrawRequest{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.Wallet{
		ID: res.Id,
	}

	return http.StatusOK, &response
}

func (u *paymentService) CreateInvoice(input dto.CreateInvoiceRequest) (int, *dto.CreateInvoiceResponse) {
	res, err := u.Client.CreateInvoice(context.Background(), &pb.CreateInvoiceRequest{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.CreateInvoiceResponse{
		InvoiceUrl: res.InvoiceUrl,
	}

	return http.StatusOK, &response
}
