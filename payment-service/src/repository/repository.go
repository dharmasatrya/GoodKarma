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
	GetWalletByUserId(ctx context.Context, userId string) (*entity.Wallet, error)
	CheckBalanceForWithdrawal(ctx context.Context, input entity.WithdrawRequest) error
}

type paymentRepository struct {
	db *mongo.Database
}

func (r *paymentRepository) GetWalletCollection() *mongo.Collection {
	return r.db.Collection("wallets")
}

func NewPaymentRepository(db *mongo.Database) *paymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) CreateWallet(input entity.Wallet) (*entity.Wallet, error) {
	walletCollection := r.GetWalletCollection()

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

	_, err2 := walletCollection.InsertOne(ctx, wallet)
	if err2 != nil {
		return nil, err2
	}

	return &wallet, nil
}

func (r *paymentRepository) GetWalletByUserId(ctx context.Context, userId string) (*entity.Wallet, error) {
	walletCollection := r.GetWalletCollection()

	var wallet entity.Wallet
	err := walletCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *paymentRepository) UpdateWalletBalance(ctx context.Context, input entity.UpdateWalleetBalanceRequest) (*entity.Wallet, error) {
	walletCollection := r.GetWalletCollection()

	var wallet entity.Wallet

	err := walletCollection.FindOne(ctx, bson.M{"user_id": input.UserID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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
	err1 := walletCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": wallet.ID},
		bson.M{"$set": bson.M{"amount": wallet.Amount}}, // Added the update operation
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedWallet)

	if err1 != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	return &updatedWallet, nil
}

func (r *paymentRepository) CheckBalanceForWithdrawal(ctx context.Context, input entity.WithdrawRequest) error {
	walletCollection := r.GetWalletCollection()

	// First find the wallet
	var wallet entity.Wallet
	err := walletCollection.FindOne(ctx, bson.M{"user_id": input.UserId}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return status.Errorf(codes.NotFound, "wallet not found")
		}
		return status.Errorf(codes.Internal, "failed to find wallet: %v", err)
	}

	// Check if there's sufficient balance
	if wallet.Amount < input.Amount {
		return status.Errorf(codes.FailedPrecondition, "insufficient balance")
	}

	return nil
}
