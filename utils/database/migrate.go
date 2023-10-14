package database

import (
	admin "restoran/features/admin/model"
	menu "restoran/features/menu/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&menu.Menu{}, &admin.Admin{})
}
