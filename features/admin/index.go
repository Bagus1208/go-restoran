package admin

import (
	"restoran/config"
	"restoran/features/admin/handler"
	"restoran/features/admin/repository"
	"restoran/features/admin/service"
	"restoran/helper"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func FeatureAdmin(db *gorm.DB, config config.Config, validate *validator.Validate) handler.AdminHandlerInterface {
	var jwt = helper.NewJWT(config.Secret)
	var generate = helper.NewGenerator()

	var adminModel = repository.NewAdminRepo(db)
	var adminService = service.NewAdminService(adminModel, generate, jwt, validate)
	var adminHandler = handler.NewAdminHandler(adminService)

	return adminHandler
}
