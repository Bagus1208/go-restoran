package model

type TransactionResponse struct {
	ID      string `json:"id"`
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
}

type TransactionInputResponse struct {
	ID          string `json:"id"`
	OrderID     uint   `json:"order_id"`
	Status      string `json:"status"`
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}
