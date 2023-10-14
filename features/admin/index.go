package admin

import (
	"restoran/config"
	"restoran/features/admin/handler"
	"restoran/features/admin/repository"
	"restoran/features/admin/service"
	"restoran/helper"

	"gorm.io/gorm"
)

func FeatureAdmin(db *gorm.DB, config config.Config) handler.AdminHandlerInterface {
	var jwt = helper.NewJWT(config.Secret)
	var generate = helper.NewGenerator()

	var adminModel = repository.NewAdminRepo(db)
	var adminService = service.NewAdminService(adminModel, generate, jwt)
	var adminHandler = handler.NewAdminHandler(adminService)

	return adminHandler
}
