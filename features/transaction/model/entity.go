package model

import "restoran/features/order/model"

type Transaction struct {
	ID      string      `gorm:"primaryKey;type:varchar(255)"`
	OrderID uint        `gorm:"type:int; not null"`
	Order   model.Order `gorm:"foreignKey:OrderID"`
	Status  string      `gorm:"type:varchar(20);default:'pending'"`
}

type Status struct {
	Transaction string
	Order       string
}

type Order struct {
	ID    uint
	Total int64
}
