package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Karma struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Amount int                `bson:"amount"`
}

type CreateKarmaRequest struct {
	UserID string `json:"user_id" bson:"user_id"`
	Amount int    `json:"amount" bson:"amount"`
}

type UpdateKarmaRequest struct {
	UserID string `json:"user_id" bson:"user_id"`
	Amount int    `json:"amount" bson:"amount"`
}

type ReferralLog struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	ReferralCode string             `json:"referral_code" bson:"referral_code"`
}

type KarmaReward struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Amount      int                `json:"amount" bson:"amount"`
	Description string             `json:"description" bson:"description"`
}

type KarmaTransaction struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	KarmaRewardID primitive.ObjectID `json:"karma_reward_id" bson:"karma_reward_id"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}

type ExchangeRewardRequest struct {
	UserID        string `json:"user_id" bson:"user_id"`
	KarmaRewardID string `json:"karma_reward_id" bson:"karma_reward_id"`
}

type CashbackDonationRequest struct {
	UserID string `json:"user_id" bson:"user_id"`
	Amount int    `json:"amount" bson:"amount"`
}
