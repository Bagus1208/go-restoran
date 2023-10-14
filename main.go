package main

import (
	"fmt"
	"restoran/config"
	"restoran/features/menu/handler"
	"restoran/features/menu/repository"
	"restoran/features/menu/service"
	"restoran/routes"
	"restoran/utils"
	"restoran/utils/database"

	"github.com/labstack/echo/v4"
)

func main() {
	var e = echo.New()
	var config = config.InitConfig()

	var db = database.InitDB(*config)
	database.Migrate(db)

	var cdn = utils.CloudinaryInstance(*config)

	var menuModel = repository.NewMenuRepo(db, cdn)
	var menuService = service.NewMenuService(menuModel)
	var menuHandler = handler.NewMenuHandler(menuService)

	routes.RouteMenu(e, menuHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Server_Port)).Error())
}
