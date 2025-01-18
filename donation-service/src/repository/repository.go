package repository

import (
	"context"
	"time"

	"github.com/dharmasatrya/goodkarma/donation-service/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type DonationRepository interface {
	CreateDonation(input entity.Donation) (*entity.Donation, error)
}

type donationRepository struct {
	db *mongo.Collection
}

func NewDonationRepository(db *mongo.Collection) *donationRepository {
	return &donationRepository{db}
}

func (r *donationRepository) CreateDonation(input entity.Donation) (*entity.Donation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	donation := entity.Donation{
		ID:           input.ID,
		UserID:       input.UserID,
		EventID:      input.EventID,
		Amount:       input.Amount,
		Status:       input.Status,
		DonationType: input.DonationType,
	}

	_, err2 := r.db.InsertOne(ctx, donation)
	if err2 != nil {
		return nil, err2
	}

	return &donation, nil
}
