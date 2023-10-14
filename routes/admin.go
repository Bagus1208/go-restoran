package routes

import (
	"restoran/features/admin/handler"

	"github.com/labstack/echo/v4"
)

func RouteAdmin(e *echo.Echo, menuHandler handler.AdminHandlerInterface) {
	admin := e.Group("/admins")
	admin.POST("", menuHandler.Insert())
	admin.POST("/login", menuHandler.Login())
}
