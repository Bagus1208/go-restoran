package model

type MenuInput struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Category    string `json:"category" form:"category" validate:"required"`
	Price       int    `json:"price" form:"price" validate:"required"`
	Description string `json:"descriptions" form:"description" validate:"required"`
	Image       string `json:"image" form:"image"`
}

type Pagination struct {
	Page     int `json:"page" form:"page" validate:"required"`
	PageSize int `json:"pageSize" form:"pageSize" validate:"required"`
}
