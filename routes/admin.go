package routes

import (
	"os"
	"restoran/features/admin/handler"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteAdmin(e *echo.Echo, menuHandler handler.AdminHandlerInterface) {
	admin := e.Group("/admins")
	admin.POST("", menuHandler.Insert())
	admin.POST("/login", menuHandler.Login())
	admin.POST("/table", menuHandler.SetNoTable(), echojwt.JWT([]byte(os.Getenv("SECRET"))))
}
