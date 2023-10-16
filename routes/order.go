package routes

import (
	"restoran/config"
	"restoran/features/order/handler"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteOrder(e *echo.Echo, orderHandler handler.OrderHandlerInterface, config config.Config) {
	e.POST("/orders", orderHandler.Insert())

	var Orders = e.Group("/orders")
	Orders.Use(echojwt.JWT([]byte(config.Secret)))
	Orders.GET("", orderHandler.GetAll())
	Orders.GET("/:id", orderHandler.GetByID())
	Orders.PUT("/:id", orderHandler.Update())
	Orders.DELETE("/:id", orderHandler.Delete())
}
