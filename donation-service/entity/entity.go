package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Donation struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	EventID      string             `json:"event_id" bson:"event_id"`
	Amount       uint32             `json:"amount" bson:"amount"`
	Status       string             `json:"status" bson:"status"`
	DonationType string             `json:"donation_type" bson:"donation_type"`
}
