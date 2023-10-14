package menu

import (
	"restoran/features/menu/handler"
	"restoran/features/menu/repository"
	"restoran/features/menu/service"

	"github.com/cloudinary/cloudinary-go"
	"gorm.io/gorm"
)

func FeatureMenu(db *gorm.DB, cdn *cloudinary.Cloudinary) handler.MenuHandlerInterface {
	var menuModel = repository.NewMenuRepo(db, cdn)
	var menuService = service.NewMenuService(menuModel)
	var menuHandler = handler.NewMenuHandler(menuService)

	return menuHandler
}
