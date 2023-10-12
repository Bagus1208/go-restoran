package database

import (
	"restoran/features/menu/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Menu{})
}
