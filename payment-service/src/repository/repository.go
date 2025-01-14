package repository

import (
	"context"
	"goodkarma-payment-service/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
	CreateWallet(input entity.Wallet) (*entity.Wallet, error)
	// UpdateWalletBalance(ctx context.Context, input entity.UpdateWalleetBalanceRequest) (*entity.Wallet, error)
	// Withdraw(userId string) (*entity.Wallet, error)
	// CreateInvoice(input entity.Invoice) (string, error)
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
