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
	CreateInvoice(token string, input dto.CreateInvoiceRequest) (int, *dto.CreateInvoiceResponse)
	UpdateWalletBalance(callbackToken string, input dto.UpdateWalletBalanceRequest) (int, *dto.Wallet)
	GetWalletByUserId(token string) (int, *dto.Wallet)
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
		log.Fatalf("error while create request %v", err)
	}

	response := dto.WithdrawResponse{
		Message: res.Message,
	}

	return http.StatusOK, &response
}

func (u *paymentService) CreateInvoice(token string, input dto.CreateInvoiceRequest) (int, *dto.CreateInvoiceResponse) {
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := u.Client.CreateInvoice(ctx, &pb.CreateInvoiceRequest{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.CreateInvoiceResponse{
		InvoiceUrl: res.InvoiceUrl,
	}

	return http.StatusOK, &response
}

func (u *paymentService) UpdateWalletBalance(callbackToken string, input dto.UpdateWalletBalanceRequest) (int, *dto.Wallet) {

	// Verify the token matches your expected token from Xendit
	expectedToken := os.Getenv("XENDIT_CALLBACK_TOKEN")
	if callbackToken != expectedToken {
		return http.StatusForbidden, nil
	}

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

func (u *paymentService) GetWalletByUserId(token string) (int, *dto.Wallet) {
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	fmt.Println(token)

	res, err := u.Client.GetWalletByUserId(ctx, &emptypb.Empty{})
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
