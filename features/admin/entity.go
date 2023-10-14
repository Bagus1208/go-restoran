package admin

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        string `gorm:"primaryKey; type:varchar(255)" json:"id"`
	Name      string `gorm:"type:varchar(255)" json:"name"`
	Email     string `gorm:"type:varchar(255); uniqueIndex" json:"email"`
	Password  string `gorm:"type:varchar(255);" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type UserCredential struct {
	Name   string
	Access map[string]any
}
