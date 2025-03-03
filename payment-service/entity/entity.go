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
	UserID string `json:"user_id" bson:"user_id"`
	Amount uint32 `json:"amount" bson:"amount"`
	Type   string `json:"type" bson:"type"`
}

type XenditInvoiceCallbackRequest struct {
	UserID     string `json:"user_id" bson:"user_id"`
	Amount     uint32 `json:"amount" bson:"amount"`
	Type       string `json:"type" bson:"type"`
	DonationID string `json:"donation_id" bson:"donation_id"`
}

type WithdrawRequest struct {
	UserId string `json:"user_id" bson:"user_id"`
	Amount uint32 `json:"amount" bson:"amount"`
}

type XenditInvoiceRequest struct {
	ExternalId  string
	Amount      int
	Description string
	Name        string
	Email       string
	Phone       string
}

type XenditDisbursementRequest struct {
	ExternalId        string
	Amount            int
	BankCode          string
	AccountHolderName string
	BankAccountNumber string
	Description       string
	Email             string
}

type XenditInvoiceResponse struct {
	InvoiceURL string `json:"invoice_url"`
}
