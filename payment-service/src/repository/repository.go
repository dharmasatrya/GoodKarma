package repository

import (
	"context"
	"fmt"
	"log"
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
	CreateKarma(ctx context.Context, payload entity.CreateKarmaRequest) (*entity.Karma, error)
	GetReferralCount(ctx context.Context, referrerID string) (uint32, error)
	CreateReferralLog(ctx context.Context, payload entity.ReferralLog) error
	UpdateKarmaAmount(ctx context.Context, payload entity.UpdateKarmaRequest) error
	GetUserByReferralCode(ctx context.Context, referralCode string) (string, error)
}

func (r *paymentRepository) GetPaymentCollection() *mongo.Collection {
	return r.db.Collection("payments")
}

func (r *paymentRepository) GetKarmaCollection() *mongo.Collection {
	return r.db.Collection("karma")
}

func (r *paymentRepository) GetReferralLogsCollection() *mongo.Collection {
	return r.db.Collection("referral_logs")
}

type paymentRepository struct {
	db *mongo.Database
}

func NewPaymentRepository(db *mongo.Database) *paymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) CreateWallet(input entity.Wallet) (*entity.Wallet, error) {
	paymentCollection := r.GetPaymentCollection()

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

	_, err2 := paymentCollection.InsertOne(ctx, wallet)
	if err2 != nil {
		return nil, err2
	}

	return &wallet, nil
}

func (r *paymentRepository) GetWalletByUserId(ctx context.Context, userId string) (*entity.Wallet, error) {
	paymentCollection := r.GetPaymentCollection()

	var wallet entity.Wallet
	err := paymentCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *paymentRepository) UpdateWalletBalance(ctx context.Context, input entity.UpdateWalleetBalanceRequest) (*entity.Wallet, error) {
	paymentCollection := r.GetPaymentCollection()

	var wallet entity.Wallet

	err := paymentCollection.FindOne(ctx, bson.M{"user_id": input.UserID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("error 62", err, input.UserID)
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
	err1 := paymentCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": wallet.ID},
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

func (r *paymentRepository) CheckBalanceForWithdrawal(ctx context.Context, input entity.WithdrawRequest) error {
	paymentCollection := r.GetPaymentCollection()

	// First find the wallet
	var wallet entity.Wallet
	err := paymentCollection.FindOne(ctx, bson.M{"user_id": input.UserId}).Decode(&wallet)
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

func (r *paymentRepository) CreateKarma(ctx context.Context, payload entity.CreateKarmaRequest) (*entity.Karma, error) {
	karmaCollection := r.GetKarmaCollection()

	userID, err := primitive.ObjectIDFromHex(payload.UserID)

	karma := entity.Karma{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		Amount: payload.Amount,
	}

	_, err = karmaCollection.InsertOne(ctx, karma)
	if err != nil {
		return nil, fmt.Errorf("failed to create karma: %v", err)
	}

	return &karma, nil
}

func (r *paymentRepository) GetReferralCount(ctx context.Context, referrerID string) (uint32, error) {
	referralCollection := r.GetReferralLogsCollection()

	// Find all referral logs for the referrer
	count, err := referralCollection.CountDocuments(ctx, bson.M{"referral_code": referrerID})

	if err != nil {
		return 0, fmt.Errorf("failed to count referral logs: %v", err)
	}

	return uint32(count), nil
}

func (r *paymentRepository) CreateReferralLog(ctx context.Context, payload entity.ReferralLog) error {
	referralCollection := r.GetReferralLogsCollection()

	referralLog := entity.ReferralLog{
		ID:           primitive.NewObjectID(),
		UserID:       payload.UserID,
		ReferralCode: payload.ReferralCode,
	}

	_, err := referralCollection.InsertOne(ctx, referralLog)
	if err != nil {
		return fmt.Errorf("failed to create referral log: %v", err)
	}

	return nil
}

func (r *paymentRepository) UpdateKarmaAmount(ctx context.Context, payload entity.UpdateKarmaRequest) error {
	karmaCollection := r.GetKarmaCollection()

	userID, err := primitive.ObjectIDFromHex(payload.UserID)

	if err != nil {
		return err
	}

	filter := bson.M{"user_id": userID}
	update := bson.M{"$inc": bson.M{"amount": payload.Amount}}

	log.Printf("Filter: %+v, Update: %+v", filter, update)

	_, err = karmaCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) GetUserByReferralCode(ctx context.Context, referralCode string) (string, error) {
	karmaCollection := r.GetKarmaCollection()

	var karma entity.Karma

	err := karmaCollection.FindOne(ctx, bson.M{"referral_code": referralCode}).Decode(&karma)

	if err != nil {
		return "", err
	}

	return karma.UserID.Hex(), nil
}
