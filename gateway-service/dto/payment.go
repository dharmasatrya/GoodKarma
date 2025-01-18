package dto

type Wallet struct {
	ID           string `json:"id" bson:"id"`
	UserID       string `json:"user_id" bson:"user_id"`
	EventID      string `json:"event_id" bson:"event_id"`
	Amount       string `json:"amount" bson:"amount"`
	DonationType string `json:"donation_type" bson:"donation_type"`
}

type CreateWalletRequest struct {
	UserID            string `json:"user_id" bson:"user_id"`
	BankAccountName   string `json:"bank_account_name" bson:"bank_account_name"`
	BankCode          string `json:"bank_code" bson:"bank_code"`
	BankAccountNumber string `json:"bank_account_number" bson:"bank_account_number"`
}

type UpdateWalletBalanceRequest struct {
	ID     string `json:"id" bson:"id"`
	Amount uint32 `json:"amount" bson:"amount"`
	Type   string `json:"type" bson:"type"` //MONEYIN or MONEYOUT
}

type WithdrawRequest struct {
	Amount uint32 `json:"amount" bson:"amount"`
}

type CreateInvoiceRequest struct {
	UserID      string `json:"user_id" bson:"user_id"`
	ExternalID  string `json:"external_id" bson:"external_id"`
	Amount      uint32 `json:"amount" bson:"amount"`
	Description string `json:"description" bson:"description"`
}

type CreateInvoiceResponse struct {
	InvoiceUrl string `json:"invoice_url" bson:"invoice_url"`
}
