package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dharmasatrya/goodkarma/payment-service/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentRepository interface {
	CreateWallet(input entity.Wallet) (*entity.Wallet, error)
	UpdateWalletBalance(ctx context.Context, input entity.UpdateWalleetBalanceRequest) (*entity.Wallet, error)
	Withdraw(ctx context.Context, input entity.WithdrawRequest) (*entity.Wallet, error)
}

type paymentRepository struct {
	db *mongo.Collection
}

func NewPaymentRepository(db *mongo.Collection) *paymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) CreateWallet(input entity.Wallet) (*entity.Wallet, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wallet := entity.Wallet{
		ID:                primitive.NewObjectID(),
		UserID:            input.UserID,
		BankAccountName:   input.BankAccountName,
		BankCode:          input.BankCode,
		BankAccountNumber: input.BankAccountNumber,
		Amount:            0,
	}

	_, err2 := r.db.InsertOne(ctx, wallet)
	if err2 != nil {
		return nil, err2
	}

	return &wallet, nil
}

func (r *paymentRepository) UpdateWalletBalance(ctx context.Context, input entity.UpdateWalleetBalanceRequest) (*entity.Wallet, error) {
	var wallet entity.Wallet

	err := r.db.FindOne(ctx, bson.M{"user_id": input.UserID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("error 62", err)
			return nil, fmt.Errorf("wallet not found")
		}
	}

	if input.Type == "money_in" {
		wallet.Amount += input.Amount
	}

	if input.Type == "money_out" {
		wallet.Amount -= input.Amount
	}

	var updatedWallet entity.Wallet
	err1 := r.db.FindOneAndUpdate(
		ctx,
		bson.M{"id": wallet.ID},
		bson.M{"$set": bson.M{"amount": wallet.Amount}}, // Added the update operation
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedWallet)

	if err1 != nil {
		fmt.Println("error 82", err1)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	return &updatedWallet, nil
}

func (r *paymentRepository) Withdraw(ctx context.Context, input entity.WithdrawRequest) (*entity.Wallet, error) {
	// First find the wallet
	var wallet entity.Wallet
	err := r.db.FindOne(ctx, bson.M{"user_id": input.UserId}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "wallet not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find wallet: %v", err)
	}

	// Check if there's sufficient balance
	if wallet.Amount < input.Amount {
		return nil, status.Errorf(codes.FailedPrecondition, "insufficient balance")
	}

	// Calculate new amount
	newAmount := wallet.Amount - input.Amount

	// Update the wallet with new amount
	var updatedWallet entity.Wallet
	err = r.db.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": input.UserId},
		bson.M{
			"$set": bson.M{
				"amount": newAmount,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedWallet)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "wallet not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update wallet: %v", err)
	}

	return &updatedWallet, nil
}
