package routes

import (
	"restoran/config"
	"restoran/features/menu/handler"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteMenu(e *echo.Echo, menuHandler handler.MenuHandlerInterface, config config.Config) {
	e.GET("/menus", menuHandler.GetData())
	e.GET("/menus/favorites", menuHandler.GetFavorite())
	e.POST("/menus/recommendations", menuHandler.RecommendationMenu())

	menus := e.Group("/menus")
	menus.Use(echojwt.JWT([]byte(config.Secret)))
	menus.POST("", menuHandler.Insert())
	menus.PUT("/:id", menuHandler.Update())
	menus.DELETE("/:id", menuHandler.Delete())
}
