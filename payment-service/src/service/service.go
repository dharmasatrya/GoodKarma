// service/borrow_service.go
package service

import (
	"context"
	"goodkarma-payment-service/entity"
	"goodkarma-payment-service/src/repository"

	pb "goodkarma-payment-service/proto"

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
