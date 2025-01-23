package repository

import (
	"context"
	"fmt"

	"github.com/dharmasatrya/goodkarma/karma-service/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KarmaRepository interface {
	CreateKarma(context.Context, entity.CreateKarmaRequest) (*entity.Karma, error)
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
