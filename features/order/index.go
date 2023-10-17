package order

import (
	"restoran/config"
	"restoran/features/order/handler"
	"restoran/features/order/repository"
	"restoran/features/order/service"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func FeatureOrder(db *gorm.DB, config config.Config, validate *validator.Validate) handler.OrderHandlerInterface {
	var orderModel = repository.NewOrderRepo(db)
	var orderService = service.NewOrderService(orderModel, validate)
	var orderHandler = handler.NewOrderHandler(orderService)

	return orderHandler
}
