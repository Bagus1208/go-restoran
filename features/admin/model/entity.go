package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        string         `gorm:"primaryKey; type:varchar(255)" json:"id"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Email     string         `gorm:"type:varchar(255); uniqueIndex" json:"email"`
	Password  string         `gorm:"type:varchar(255);" json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type UserCredential struct {
	Name   string         `json:"name"`
	Access map[string]any `json:"access"`
}
