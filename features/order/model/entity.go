package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NoTable   int            `json:"no_table"`
	Orders    []OrderDetail  `gorm:"foreignKey:OrderID" json:"orders"`
	Total     int            `json:"total"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type OrderDetail struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   int            `json:"order_id"`
	MenuName  string         `json:"menu_name"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
