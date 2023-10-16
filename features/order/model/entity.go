package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	NoTable int           `json:"no_table"`
	Orders  []OrderDetail `gorm:"foreignKey:OrderID" json:"orders"`
	Total   int           `json:"total"`
}

type OrderDetail struct {
	gorm.Model
	OrderID  int    `json:"order_id"`
	MenuName string `json:"menu_name"`
	Quantity int    `json:"quantity"`
}

type MenuPrice struct {
	Name  string
	Price int
}
