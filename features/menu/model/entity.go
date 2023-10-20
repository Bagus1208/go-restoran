package model

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID          uint   `gorm:"primaryKey; type:int"`
	Name        string `gorm:"type:varchar(255)"`
	Category    string `gorm:"type:varchar(255)"`
	Price       int    `gorm:"type:int"`
	Description string `gorm:"type:text"`
	Image       string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
