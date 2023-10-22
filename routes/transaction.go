package routes

import (
	"restoran/config"
	"restoran/features/transaction/handler"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteTransaction(e *echo.Echo, transactionHandler handler.TransactionHandlerInterface, config config.Config) {
	e.POST("/transactions/notifications", transactionHandler.Notifications())

	var Transactions = e.Group("/transactions")
	Transactions.Use(echojwt.JWT([]byte(config.Secret)))
	Transactions.POST("", transactionHandler.Insert())
	Transactions.GET("", transactionHandler.GetAll())
	Transactions.GET("/:id", transactionHandler.GetByID())
	Transactions.DELETE("/:id", transactionHandler.Delete())
}
