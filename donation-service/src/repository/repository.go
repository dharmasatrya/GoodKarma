package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dharmasatrya/goodkarma/donation-service/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DonationRepository interface {
	CreateDonation(input entity.Donation) (*entity.Donation, error)
	UpdateDonationStatus(input entity.Donation) (*entity.Donation, error)
	GetDonationsByUserId(userId string) ([]entity.Donation, error)
	GetDonationsByEventId(eventId string) ([]entity.Donation, error)
	CreateDonationWithSession(ctx context.Context, donation entity.Donation) (*entity.Donation, error)
	StartSession() (mongo.Session, error)
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

func (r *donationRepository) UpdateDonationStatus(input entity.Donation) (*entity.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": input.ID}
	update := bson.M{"$set": bson.M{"status": input.Status}}

	// Update the document
	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Fetch the updated document
	var updatedDonation entity.Donation
	err = r.db.FindOne(ctx, filter).Decode(&updatedDonation)
	if err != nil {
		return nil, err
	}

	return &updatedDonation, nil
}

func (r *donationRepository) GetDonationsByUserId(userId string) ([]entity.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var donations []entity.Donation
	if err = cursor.All(ctx, &donations); err != nil {
		return nil, err
	}

	return donations, nil
}

func (r *donationRepository) GetDonationsByEventId(eventId string) ([]entity.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"event_id": eventId}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var donations []entity.Donation
	if err = cursor.All(ctx, &donations); err != nil {
		return nil, err
	}

	return donations, nil
}

func (r *donationRepository) CreateDonationWithSession(ctx context.Context, donation entity.Donation) (*entity.Donation, error) {
	_, err := r.db.InsertOne(ctx, donation)
	if err != nil {
		return nil, fmt.Errorf("failed to insert donation: %w", err)
	}
	return &donation, nil
}

func (r *donationRepository) StartSession() (mongo.Session, error) {
	return r.db.Database().Client().StartSession()
}
