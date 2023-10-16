package main

import (
	"fmt"
	"restoran/config"
	"restoran/features/admin"
	"restoran/features/menu"
	"restoran/features/order"
	"restoran/helper"
	"restoran/routes"
	"restoran/utils"
	"restoran/utils/database"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	var e = echo.New()
	var config = config.InitConfig()
	var validate = validator.New()

	var db = database.InitDB(*config)
	database.Migrate(db)

	var cdn = utils.CloudinaryInstance(*config)

	var menuHandler = menu.FeatureMenu(db, cdn, validate)
	var adminHandler = admin.FeatureAdmin(db, *config, validate)
	var orderHandler = order.FeatureOrder(db)

	helper.LogMiddlewares(e)

	routes.RouteMenu(e, menuHandler, *config)
	routes.RouteAdmin(e, adminHandler)
	routes.RouteOrder(e, orderHandler, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Server_Port)).Error())
}
