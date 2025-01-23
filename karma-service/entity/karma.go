package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Karma struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Amount uint32             `bson:"amount"`
}

type CreateKarmaRequest struct {
	UserID string `json:"user_id" bson:"user_id"`
	Amount uint32 `json:"amount" bson:"amount"`
}

type UpdateKarmaRequest struct {
	UserID string `json:"user_id" bson:"user_id"`
	Amount uint32 `json:"amount" bson:"amount"`
}

type ReferralLog struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	ReferralCode string             `json:"referral_code" bson:"referral_code"`
}
