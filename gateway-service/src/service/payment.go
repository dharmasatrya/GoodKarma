package service

import (
	"context"
	"fmt"
	"gateway-service/dto"
	"log"
	"net/http"
	"os"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PaymentService interface {
	Withdraw(token string, input dto.WithdrawRequest) (int, *dto.WithdrawResponse)
	GetWalletByUserId(token string) (int, *dto.Wallet)
	UpdateInvoiceWalletBalance(callbackToken string, input dto.UpdateInvoiceBalanceRequest) (int, *dto.Donation)
	UpdateDisbursementWalletBalance(callbackToken string, input dto.XenditDisbursementCallbackRequest) (int, *dto.Wallet)
}

type paymentService struct {
	Client pb.PaymentServiceClient
}

func NewPaymentService(paymentClient pb.PaymentServiceClient) *paymentService {
	return &paymentService{paymentClient}
}

func (u *paymentService) Withdraw(token string, input dto.WithdrawRequest) (int, *dto.WithdrawResponse) {
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := u.Client.Withdraw(ctx, &pb.WithdrawRequest{
		Amount: input.Amount,
	})
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	response := dto.WithdrawResponse{
		Message: res.Message,
	}

	return http.StatusOK, &response
}

func (u *paymentService) UpdateDisbursementWalletBalance(callbackToken string, input dto.XenditDisbursementCallbackRequest) (int, *dto.Wallet) {

	// Verify the token matches your expected token from Xendit
	expectedToken := os.Getenv("XENDIT_CALLBACK_TOKEN")
	if callbackToken != expectedToken {
		return http.StatusForbidden, nil
	}

	res, err := u.Client.XenditDisbursementCallback(context.Background(), &pb.XenditDisbursementCallbackRequest{
		ExternalId: input.ExternalID,
		Amount:     input.Amount,
		Type:       input.Type,
	})

	if err != nil {
		return http.StatusInternalServerError, nil
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

func (u *paymentService) UpdateInvoiceWalletBalance(callbackToken string, input dto.UpdateInvoiceBalanceRequest) (int, *dto.Donation) {

	expectedToken := os.Getenv("XENDIT_CALLBACK_TOKEN")
	if callbackToken != expectedToken {
		return http.StatusForbidden, nil
	}

	res, err := u.Client.XenditInvoiceCallback(context.Background(), &pb.XenditInvoiceCallbackRequest{
		Amount:     input.Amount,
		Type:       input.Type,
		DonationId: input.DonationID,
	})
	if err != nil {
		log.Println("XenditInvoiceCallback err %v:", err)
		return http.StatusInternalServerError, nil
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

func (u *paymentService) GetWalletByUserId(token string) (int, *dto.Wallet) {
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	fmt.Println(token)

	res, err := u.Client.GetWalletByUserId(ctx, &emptypb.Empty{})
	if err != nil {
		return http.StatusInternalServerError, nil
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
