package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UserID            string             `bson:"user_id"`
	BankAccountName   string             `bson:"bank_account_name"`
	BankCode          string             `bson:"bank_code"`
	BankAccountNumber string             `bson:"bank_account_number"`
	Amount            uint32             `bson:"amount"`
}

type Invoice struct {
	UserID      uint32 `bson:"user_id"`
	ExternalID  string `bson:"external_id"`
	Amount      uint32 `bson:"amount"`
	Description string `bson:"description"`
}

type UpdateWalleetBalanceRequest struct {
	UserID     uint32 `bson:"user_id"`
	ExternalID string `bson:"external_id"`
	Amount     uint32 `bson:"amount"`
	Type       string `bson:"type"`
}
