package model

type MenuInput struct {
	Name        string `json:"name" form:"name"`
	Category    string `json:"category" form:"category"`
	Price       int    `json:"price" form:"price"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
}
