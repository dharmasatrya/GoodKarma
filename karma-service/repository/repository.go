package repository

import (
	"context"
	"fmt"

	"github.com/dharmasatrya/goodkarma/karma-service/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KarmaRepository interface {
	CreateKarma(context.Context, entity.CreateKarmaRequest) (*entity.Karma, error)
	GetReferralCount(context.Context, string) (uint32, error)
	CreateReferralLog(context.Context, entity.ReferralLog) error
}

type karmaRepository struct {
	db *mongo.Database
}

func NewKarmaRepository(db *mongo.Database) KarmaRepository {
	return &karmaRepository{db}
}

func (r *karmaRepository) GetKarmaCollection() *mongo.Collection {
	return r.db.Collection("karma")
}

func (r *karmaRepository) GetReferralLogsCollection() *mongo.Collection {
	return r.db.Collection("referral_logs")
}

func (r *karmaRepository) CreateKarma(ctx context.Context, payload entity.CreateKarmaRequest) (*entity.Karma, error) {
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

func (r *karmaRepository) GetReferralCount(ctx context.Context, referrerID string) (uint32, error) {
	referralCollection := r.GetReferralLogsCollection()

	// Find all referral logs for the referrer
	count, err := referralCollection.CountDocuments(ctx, bson.M{"referral_code": referrerID})

	if err != nil {
		return 0, fmt.Errorf("failed to count referral logs: %v", err)
	}

	return uint32(count), nil
}

func (r *karmaRepository) CreateReferralLog(ctx context.Context, payload entity.ReferralLog) error {
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
