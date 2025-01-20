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
	UpdateWalletBalance(input dto.UpdateWalletBalanceRequest) (int, *dto.Wallet)
}

type paymentService struct {
	Client pb.PaymentServiceClient
}

func NewPaymentService(paymentClient pb.PaymentServiceClient) *paymentService {
	return &paymentService{paymentClient}
}

func (u *paymentService) Withdraw(input dto.WithdrawRequest) (int, *dto.WithdrawResponse) {
	res, err := u.Client.Withdraw(context.Background(), &pb.WithdrawRequest{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.WithdrawResponse{
		Message: res.Message,
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

func (u *paymentService) UpdateWalletBalance(input dto.UpdateWalletBalanceRequest) (int, *dto.Wallet) {

	res, err := u.Client.UpdateWalletBalance(context.Background(), &pb.UpdateWalletBalanceRequest{
		Amount: input.Amount,
		Type:   input.Type,
	})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.Wallet{
		ID:                res.Id,
		UserID:            res.UserId,
		BankAccountName:   res.BankAccountName,
		BankCode:          res.BankCode,
		BankAccountNumber: res.BankAccountNumber,
		Amount:            res.Amount,
	}

	return http.StatusOK, &response
}
