package model

type OrderInput struct {
	Orders  []OrderDetail `json:"orders" form:"orders" validate:"required"`
	NoTable int           `json:"no_table" form:"no_table" validate:"required"`
}

type Pagination struct {
	Page     int
	PageSize int
}
