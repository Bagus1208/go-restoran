package model

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name        string         `gorm:"type:varchar(255)" json:"name"`
	Category    string         `gorm:"type:varchar(255)" json:"category"`
	Price       int            `gorm:"type:int" json:"price"`
	Description string         `gorm:"type:text" json:"description"`
	Image       string         `json:"image"`
}
