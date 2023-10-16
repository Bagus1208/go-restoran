package menu

import (
	"restoran/features/menu/handler"
	"restoran/features/menu/repository"
	"restoran/features/menu/service"

	"github.com/cloudinary/cloudinary-go"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func FeatureMenu(db *gorm.DB, cdn *cloudinary.Cloudinary, validate *validator.Validate) handler.MenuHandlerInterface {
	var menuModel = repository.NewMenuRepo(db, cdn)
	var menuService = service.NewMenuService(menuModel, validate)
	var menuHandler = handler.NewMenuHandler(menuService)

	return menuHandler
}
