package database

import (
	admin "restoran/features/admin/model"
	menu "restoran/features/menu/model"
	order "restoran/features/order/model"
	transaction "restoran/features/transaction/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&menu.Menu{}, &admin.Admin{}, &order.Order{}, &order.OrderDetail{}, &transaction.Transaction{})
}
