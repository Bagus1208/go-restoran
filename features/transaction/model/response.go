package model

type TransactionResponse struct {
	ID            uint   `json:"id"`
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	Status        string `json:"status"`
}

type TransactionInputResponse struct {
	ID          uint   `json:"id"`
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}
