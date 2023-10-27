package model

type TransactionInput struct {
	OrderID       uint   `json:"order_id" form:"order_id" validate:"required"`
	PaymentMethod string `json:"payment_method" form:"payment_method"`
}

type QueryParam struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}
