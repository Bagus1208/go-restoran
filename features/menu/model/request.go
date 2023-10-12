package model

type MenuInput struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}
