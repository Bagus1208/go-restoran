package model

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Category    string `gorm:"type:varchar(255)" json:"category"`
	Price       int    `gorm:"type:int" json:"price"`
	Description string `gorm:"type:text" json:"description"`
}
