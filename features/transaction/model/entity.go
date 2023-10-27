package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID            uint   `gorm:"primaryKey;type: int"`
	OrderID       string `gorm:"type:varchar(255); not null"`
	PaymentMethod string `gorm:"type:varchar(30)"`
	Status        string `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Status struct {
	Status        string
	Order         string
	PaymentMethod string
}

type Order struct {
	ID    uint
	Total int64
}
