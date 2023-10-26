package model

type OrderInput struct {
	Orders      []OrderDetail `json:"orders" form:"orders" validate:"required"`
	TableNumber int           `json:"table_number" form:"table_number" validate:"required"`
}

type Pagination struct {
	Page     int
	PageSize int
}
