package model

type MenuInput struct {
	Name        string `form:"name" validate:"required"`
	Category    string `form:"category" validate:"required"`
	Price       int    `form:"price" validate:"required"`
	Description string `form:"description" validate:"required"`
	Image       string `form:"image"`
}

type Pagination struct {
	Page     int
	PageSize int
}
