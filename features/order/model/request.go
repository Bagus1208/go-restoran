package model

type OrderInput struct {
	Orders []OrderDetail `json:"orders" form:"orders"`
}
