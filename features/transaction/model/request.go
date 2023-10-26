package model

type TransactionInput struct {
	OrderID uint `json:"order_id" form:"order_id" validate:"required"`
}

type QueryParam struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}
