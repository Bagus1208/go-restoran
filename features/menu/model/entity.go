package model

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255)" json:"name" form:"name"`
	Category    string `gorm:"type:varchar(255)" json:"category" form:"category"`
	Price       int    `gorm:"type:int" json:"price" form:"price"`
	Description string `gorm:"type:text" json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
}
