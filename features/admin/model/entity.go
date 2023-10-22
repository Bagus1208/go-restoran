package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        string `gorm:"primaryKey; type:varchar(255)"`
	Name      string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255); uniqueIndex"`
	Password  string `gorm:"type:varchar(255);"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
