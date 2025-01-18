package entity

type UserRegistData struct {
	Email string `json:"email"`
	Link  string `json:"link"`
}

type InvoiceData struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	InvoiceID   string `json:"invoice_id"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Ammount     string `json:"ammount"`
}

type GoodsData struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Alamat  string `json:"alamat"`
	Ammount string `json:"ammount"`
}
