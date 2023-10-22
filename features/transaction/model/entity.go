package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID uint `gorm:"primaryKey;type: int"`
	// Order     model.Order `gorm:"foreignKey:ID"`
	OrderID   string `gorm:"type:varchar(255); not null"`
	Status    string `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Status struct {
	Transaction string
	Order       string
}

type Order struct {
	ID    uint
	Total int64
}
