package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Donation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       string             `bson:"user_id"`
	EventID      string             `bson:"event_id"`
	Amount       string             `bson:"amount"`
	Status       string             `bson:"status"`
	DonationType string             `bson:"donation_type"`
}
