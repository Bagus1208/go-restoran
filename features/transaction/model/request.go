package model

type TransactionInput struct {
	OrderID uint `json:"order_id" form:"order_id"`
}
