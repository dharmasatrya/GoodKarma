// service/borrow_service.go
package service

import (
	"context"

	"github.com/dharmasatrya/goodkarma/payment-service/entity"
	"github.com/dharmasatrya/goodkarma/payment-service/external"
	"github.com/dharmasatrya/goodkarma/payment-service/src/repository"

	pb "github.com/dharmasatrya/goodkarma/payment-service/proto"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	paymentRepository repository.PaymentRepository
}

// var jwtSecret = []byte("secret")

func NewPaymentService(paymentRepository repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		paymentRepository: paymentRepository,
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

func (s *PaymentService) UpdateWalletBalance(ctx context.Context, req *pb.UpdateWalletBalanceRequest) (*pb.UpdateWalleetBalanceResponse, error) {

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

	balanceShift := entity.UpdateWalleetBalanceRequest{
		UserID: userID,
		Amount: req.Amount,
		Type:   req.Type,
	}

	res, err := s.paymentRepository.UpdateWalletBalance(ctx, balanceShift)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
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

	invoice := entity.XenditInvoiceRequest{
		ExternalId:  "111",
		Amount:      int(req.Amount),
		Description: req.Description,
		FirstName:   "a",
		LastName:    "a",
		Email:       "a",
		Phone:       "",
	}

	_, err := external.CreateXenditInvoice(invoice)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	return &pb.CreateInvoiceResponse{
		InvoiceUrl: "url",
	}, nil
}
