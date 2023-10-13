package repository

import (
	"restoran/config"
	"restoran/features/menu/model"

	"gorm.io/gorm"
)

type MenuRepositoryInterface interface {
	Insert(newData *model.Menu) (*model.Menu, error)
	GetAll() ([]model.Menu, error)
	GetCategory(category string) ([]model.Menu, error)
	Update(id int, updateData *model.Menu) (*model.Menu, error)
	Delete(id int) error
}

type menuRepo struct {
	db     *gorm.DB
	config config.Config
}

func NewMenuRepo(db *gorm.DB) MenuRepositoryInterface {
	return &menuRepo{
		db: db,
	}
}

func (repository *menuRepo) Insert(newData *model.Menu) (*model.Menu, error) {
	result := repository.db.Create(newData)
	if result.Error != nil {
		return nil, result.Error
	}

	return newData, nil
}

func (repository *menuRepo) GetAll() ([]model.Menu, error) {
	var menus []model.Menu
	result := repository.db.Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}

func (repository *menuRepo) GetCategory(category string) ([]model.Menu, error) {
	var menus []model.Menu
	result := repository.db.Where("category = ?", category).Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}

func (repository *menuRepo) Update(id int, updateData *model.Menu) (*model.Menu, error) {
	result := repository.db.Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedUser = new(model.Menu)
	if err := repository.db.Where("id = ?", id).First(updatedUser).Error; err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (repository *menuRepo) Delete(id int) error {
	var deleteMenu model.Menu
	deleteMenu.ID = uint(id)
	result := repository.db.Delete(&deleteMenu)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
