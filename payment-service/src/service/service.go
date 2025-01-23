// service/borrow_service.go
package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	donationpb "github.com/dharmasatrya/goodkarma/donation-service/proto"
	eventpb "github.com/dharmasatrya/goodkarma/event-service/proto"
	"github.com/dharmasatrya/goodkarma/payment-service/client"
	"github.com/dharmasatrya/goodkarma/payment-service/entity"
	"github.com/dharmasatrya/goodkarma/payment-service/external"
	"github.com/dharmasatrya/goodkarma/payment-service/src/repository"
	"github.com/dharmasatrya/goodkarma/user-service/proto"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	paymentRepository repository.PaymentRepository
	userClient        *client.UserServiceClient
	donationClient    *client.DonationServiceClient
	eventClient       *client.EventServiceClient
	messageBroker     MessageBroker
}

// var jwtSecret = []byte("secret")

func NewPaymentService(
	paymentRepository repository.PaymentRepository,
	userClient *client.UserServiceClient,
	donationClient *client.DonationServiceClient,
	eventClient *client.EventServiceClient,
	messageBroker MessageBroker,
) *PaymentService {
	return &PaymentService{
		paymentRepository: paymentRepository,
		userClient:        userClient,
		donationClient:    donationClient,
		eventClient:       eventClient,
		messageBroker:     messageBroker,
	}
}

func (s *PaymentService) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	wallet := entity.Wallet{
		ID:                primitive.NewObjectID(),
		UserID:            req.UserId,
		BankAccountName:   req.BankAccountName,
		BankCode:          req.BankCode,
		BankAccountNumber: req.BankAccountNumber,
		Amount:            0,
	}

	res, err := s.paymentRepository.CreateWallet(wallet)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	return &pb.CreateWalletResponse{
		Id:                res.ID.Hex(),
		UserId:            res.UserID,
		BankAccountName:   res.BankAccountName,
		BankCode:          res.BankCode,
		BankAccountNumber: res.BankAccountNumber,
		Amount:            res.Amount,
	}, nil
}

func (s *PaymentService) GetWalletByUserId(ctx context.Context, req *emptypb.Empty) (*pb.GetWalletResponse, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Get claims from context that was set in auth middleware
	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Extract user_id from claims and verify it matches the request
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	wallet, err := s.paymentRepository.GetWalletByUserId(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wallet not found")
	}

	return &pb.GetWalletResponse{
		Id:                wallet.ID.Hex(),
		UserId:            wallet.UserID,
		BankAccountName:   wallet.BankAccountName,
		BankCode:          wallet.BankCode,
		BankAccountNumber: wallet.BankAccountNumber,
		Amount:            wallet.Amount,
	}, nil
}

func (s *PaymentService) UpdateWalletBalance(ctx context.Context, req *pb.UpdateWalletBalanceRequest) (*pb.UpdateWalleetBalanceResponse, error) {

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Get claims from context that was set in auth middleware
	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		fmt.Printf("Claims type: %T\n", ctx.Value("claims"))
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	balanceShift := entity.UpdateWalleetBalanceRequest{
		UserID: userID,
		Amount: req.Amount,
		Type:   req.Type,
	}

	res, err := s.paymentRepository.UpdateWalletBalance(ctx, balanceShift)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating balance")
	}

	return &pb.UpdateWalleetBalanceResponse{
		Id:                res.ID.Hex(),
		UserId:            res.UserID,
		BankAccountName:   res.BankAccountName,
		BankCode:          res.BankCode,
		BankAccountNumber: res.BankAccountNumber,
		Amount:            res.Amount,
	}, nil
}

func (s *PaymentService) CreateInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pb.CreateInvoiceResponse, error) {

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Get claims from context that was set in auth middleware
	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		fmt.Printf("Claims type: %T\n", ctx.Value("claims"))
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	userDetail, errUser := s.userClient.Client.GetUserById(ctx, &proto.GetUserByIdRequest{Id: userID})
	if errUser != nil {
		return nil, status.Errorf(codes.Internal, "cant get user detail")
	}

	invoice := entity.XenditInvoiceRequest{
		ExternalId:  req.ExternalId,
		Amount:      int(req.Amount),
		Description: req.Description,
		Name:        userDetail.FullName,
		Email:       "dharmasatrya10@gmail.com",
		Phone:       "081299640904",
	}

	res, err := external.CreateXenditInvoice(invoice)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	return &pb.CreateInvoiceResponse{
		InvoiceUrl: res.InvoiceURL,
	}, nil
}

func (s *PaymentService) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {

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

	wallet, errWallet := s.paymentRepository.GetWalletByUserId(ctx, userID)
	if errWallet != nil {
		fmt.Println(errWallet)
		return nil, status.Errorf(codes.InvalidArgument, "Cant find wallet")
	}

	errBalance := s.paymentRepository.CheckBalanceForWithdrawal(ctx, entity.WithdrawRequest{UserId: userID, Amount: req.Amount})
	if errBalance != nil {
		return nil, status.Errorf(codes.InvalidArgument, "insufficient balance")
	}

	userDetail, errUser := s.userClient.Client.GetUserById(ctx, &proto.GetUserByIdRequest{Id: userID})
	if errUser != nil {
		return nil, status.Errorf(codes.Internal, "cant get user detail")
	}

	disbursement := entity.XenditDisbursementRequest{
		ExternalId:        userDetail.Id,
		Amount:            int(req.Amount),
		BankCode:          wallet.BankCode,
		AccountHolderName: wallet.BankAccountName,
		Description:       "withdraw funds from GoodKarma",
		BankAccountNumber: wallet.BankAccountNumber,
		Email:             userDetail.Email,
	}

	err := external.CreateXenditDisbursement(disbursement)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	return &pb.WithdrawResponse{
		Message: "Disbursement created, balance will be deducted once the disbursement has been completed",
	}, nil
}

func (s *PaymentService) XenditInvoiceCallback(ctx context.Context, req *pb.XenditInvoiceCallbackRequest) (*pb.Donation, error) {

	donation, err := s.donationClient.Client.UpdateDonationStatusXendit(ctx, &donationpb.UpdateDonationStatusRequest{
		Id:     req.DonationId,
		Status: "COMPLETED",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching donation")
	}

	eventIdToInt, err := strconv.Atoi(donation.EventId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error converting event id to int")
	}
	eventID := uint32(eventIdToInt)

	event, eventErr := s.eventClient.Client.GetEventById(ctx, &eventpb.Id{Id: eventID})
	if eventErr != nil {
		return nil, status.Errorf(codes.Internal, "error fetching event")
	}

	charged, err := s.ChargeFees(ctx, &pb.ChargeFeesRequest{
		UserId: event.UserId,
		Amount: req.Amount,
		Type:   "uang",
	})

	fmt.Println(charged.AmountAfterFees, "<<<<<<<<<<<<")

	balanceShift := entity.UpdateWalleetBalanceRequest{
		UserID: event.UserId,
		Amount: charged.AmountAfterFees,
		Type:   "money_in",
	}

	_, err1 := s.paymentRepository.UpdateWalletBalance(ctx, balanceShift)
	if err1 != nil {
		return nil, status.Errorf(codes.Internal, "error updating balance")
	}

	return &pb.Donation{
		Id:           donation.Id,
		UserId:       donation.UserId,
		EventId:      donation.UserId,
		Amount:       donation.Amount,
		Status:       donation.Status,
		DonationType: donation.DonationType,
	}, nil
}

func (s *PaymentService) XenditDisbursementCallback(ctx context.Context, req *pb.XenditDisbursementCallbackRequest) (*pb.UpdateWalleetBalanceResponse, error) {

	balanceShift := entity.UpdateWalleetBalanceRequest{
		UserID: req.ExternalId,
		Amount: req.Amount,
		Type:   req.Type,
	}

	res, err := s.paymentRepository.UpdateWalletBalance(ctx, balanceShift)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating balance")
	}

	return &pb.UpdateWalleetBalanceResponse{
		Id:                res.ID.Hex(),
		UserId:            res.UserID,
		BankAccountName:   res.BankAccountName,
		BankCode:          res.BankCode,
		BankAccountNumber: res.BankAccountNumber,
		Amount:            res.Amount,
	}, nil
}

func (s *PaymentService) ChargeFees(ctx context.Context, req *pb.ChargeFeesRequest) (*pb.ChargeFeesResponse, error) {

	balanceShift := entity.UpdateWalleetBalanceRequest{
		UserID: os.Getenv("MASTER_ACCOUNT_ID"),
		Amount: 2000,
		Type:   "money_in",
	}

	_, err := s.paymentRepository.UpdateWalletBalance(ctx, balanceShift)
	if err != nil {
		log.Println("error updating balance to master account")
		return nil, status.Errorf(codes.Internal, "error updating balance to master account")
	}

	if req.Type != "uang" {
		balanceShift := entity.UpdateWalleetBalanceRequest{
			UserID: req.UserId,
			Amount: 2000,
			Type:   "money_out",
		}

		_, err := s.paymentRepository.UpdateWalletBalance(ctx, balanceShift)
		if err != nil {
			log.Println("error updating balance to master account")
			return nil, status.Errorf(codes.Internal, "error updating balance to master account")
		}
	}

	return &pb.ChargeFeesResponse{
		UserId:          req.UserId,
		AmountAfterFees: req.Amount - 2000,
	}, nil
}
