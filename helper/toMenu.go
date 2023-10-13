package helper

import "restoran/features/menu/model"

func ToMenu(data model.MenuInput) *model.Menu {
	return &model.Menu{
		Name:        data.Name,
		Category:    data.Category,
		Price:       data.Price,
		Description: data.Description,
		Image:       data.Image,
	}
}