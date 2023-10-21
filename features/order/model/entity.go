package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint          `gorm:"primaryKey; type:int"`
	NoTable   int           `gorm:"type:int"`
	Orders    []OrderDetail `gorm:"foreignKey:OrderID"`
	Total     int           `gorm:"type:int"`
	Status    string        `gorm:"type:varchar(10);default:'unpaid'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type OrderDetail struct {
	ID      uint `gorm:"primaryKey; type:int"`
	OrderID int  `gorm:"type:int"`
	MenuID  int  `gorm:"type:int" json:"menu_id"`
	// Menu      model.Menu `gorm:"foreignKey:MenuID"`
	Quantity  int `gorm:"type:int" json:"quantity"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
