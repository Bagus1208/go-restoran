package model

import (
	"restoran/features/transaction/model"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          uint              `gorm:"primaryKey; type:int"`
	Transaction model.Transaction `gorm:"foreignKey:ID"`
	NoTable     int               `gorm:"type:int"`
	Orders      []OrderDetail     `gorm:"foreignKey:OrderID"`
	Total       int               `gorm:"type:int"`
	Status      string            `gorm:"type:varchar(10);default:'unpaid'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type OrderDetail struct {
	ID        uint `gorm:"primaryKey; type:int"`
	OrderID   int  `gorm:"type:int"`
	MenuID    int  `gorm:"type:int" json:"menu_id"`
	Quantity  int  `gorm:"type:int" json:"quantity"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
