package model

import (
	"restoran/features/order/model"
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID          uint                `gorm:"primaryKey; type:int"`
	Orders      []model.OrderDetail `gorm:"foreignKey:MenuID"`
	Name        string              `gorm:"type:varchar(255)"`
	Category    string              `gorm:"type:varchar(255)"`
	Price       int                 `gorm:"type:int"`
	Description string              `gorm:"type:text"`
	Image       string              `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
