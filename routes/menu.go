package routes

import (
	"restoran/features/menu/handler"

	"github.com/labstack/echo/v4"
)

func RouteMenu(e *echo.Echo, menuHandler handler.MenuHandlerInterface) {
	menus := e.Group("/menus")
	menus.POST("", menuHandler.Insert())
	menus.GET("", menuHandler.GetAll())
	menus.GET("/:category", menuHandler.GetCategory())
	menus.PUT("/:id", menuHandler.Update())
	menus.DELETE("/:id", menuHandler.Delete())
}
