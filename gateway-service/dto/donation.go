package dto

type Donation struct {
	ID           string `json:"id" bson:"id"`
	UserID       string `json:"user_id" bson:"user_id"`
	EventID      string `json:"event_id" bson:"event_id"`
	Amount       uint32 `json:"amount" bson:"amount"`
	Status       string `json:"status" bson:"status"`
	DonationType string `json:"donation_type" bson:"donation_type"`
}

type CreateDonationRequest struct {
	EventID      string `json:"event_id" bson:"event_id"`
	Amount       uint32 `json:"amount" bson:"amount"`
	DonationType string `json:"donation_type" bson:"donation_type"`
}

type UpdateDonationStatusRequest struct {
	ID     string `json:"id" bson:"id"`
	Status string `json:"status" bson:"status"`
}
