package routes

import (
	"restoran/config"
	"restoran/features/admin/handler"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteAdmin(e *echo.Echo, adminHandler handler.AdminHandlerInterface, config config.Config) {
	admin := e.Group("/admins")
	admin.POST("", adminHandler.Insert())
	admin.POST("/login", adminHandler.Login())
	admin.POST("/table", adminHandler.SetNoTable(), echojwt.JWT([]byte(config.Secret)))
}
