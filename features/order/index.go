package order

import (
	"restoran/features/order/handler"
	"restoran/features/order/repository"
	"restoran/features/order/service"

	"gorm.io/gorm"
)

func FeatureOrder(db *gorm.DB) handler.OrderHandlerInterface {
	var orderModel = repository.NewOrderRepo(db)
	var orderService = service.NewOrderService(orderModel)
	var orderHandler = handler.NewOrderHandler(orderService)

	return orderHandler
}
