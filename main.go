package main

import (
	"fmt"
	"restoran/config"
	adminHandler "restoran/features/admin/handler"
	adminRepository "restoran/features/admin/repository"
	adminService "restoran/features/admin/service"

	menuHandler "restoran/features/menu/handler"
	menuRepository "restoran/features/menu/repository"
	menuService "restoran/features/menu/service"

	orderHandler "restoran/features/order/handler"
	orderRepository "restoran/features/order/repository"
	orderService "restoran/features/order/service"

	transactionHandler "restoran/features/transaction/handler"
	transactionRepository "restoran/features/transaction/repository"
	transactionService "restoran/features/transaction/service"

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
	var jwt = helper.NewJWT(config.Secret)

	var db = database.InitDB(*config)
	database.Migrate(db)

	var cdn = utils.CloudinaryInstance(*config)

	var snapClient = utils.MidtransSnapClient(*config)
	var coreAPIClient = utils.MidtransCoreAPIClient(*config)

	var adminRepo = adminRepository.NewAdminRepo(db)
	var adminService = adminService.NewAdminService(adminRepo, jwt, validate)
	var adminHandler = adminHandler.NewAdminHandler(adminService)

	var menuRepo = menuRepository.NewMenuRepo(db, cdn, *config)
	var menuService = menuService.NewMenuService(menuRepo, validate)
	var menuHandler = menuHandler.NewMenuHandler(menuService)

	var orderRepo = orderRepository.NewOrderRepo(db)
	var orderService = orderService.NewOrderService(orderRepo, validate, jwt)
	var orderHandler = orderHandler.NewOrderHandler(orderService)

	var transactionRepo = transactionRepository.NewTransactionRepo(db, snapClient, coreAPIClient)
	var transactionService = transactionService.NewTransactionService(transactionRepo, validate)
	var transactionHandler = transactionHandler.NewTransactionHandler(transactionService)

	helper.LogMiddlewares(e)

	routes.RouteMenu(e, menuHandler, *config)
	routes.RouteAdmin(e, adminHandler)
	routes.RouteOrder(e, orderHandler, *config)
	routes.RouteTransaction(e, transactionHandler, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Server_Port)).Error())
}
