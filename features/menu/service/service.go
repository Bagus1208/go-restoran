package service

import (
	"errors"
	"restoran/features/menu/model"
	"restoran/features/menu/repository"
)

type MenuServiceInterface interface {
	Insert(newData model.MenuInput) (*model.Menu, error)
	GetAll() ([]model.Menu, error)
	GetCategory(category string) ([]model.Menu, error)
	Update(id int, updateData model.MenuInput) (*model.Menu, error)
	Delete(id int) error
}

type menuService struct {
	repository repository.MenuRepositoryInterface
}

func NewMenuService(repo repository.MenuRepositoryInterface) MenuServiceInterface {
	return &menuService{
		repository: repo,
	}
}

func (service *menuService) Insert(newData model.MenuInput) (*model.Menu, error) {
	var newUser = new(model.Menu)
	newUser.Name = newData.Name
	newUser.Category = newData.Category
	newUser.Price = int(newData.Price)
	newUser.Description = newData.Description

	result, err := service.repository.Insert(newUser)
	if err != nil {
		return nil, errors.New("error inserting")
	}

	return result, nil
}

func (service *menuService) GetAll() ([]model.Menu, error) {
	result, err := service.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *menuService) GetCategory(category string) ([]model.Menu, error) {
	result, err := service.repository.GetCategory(category)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *menuService) Update(id int, updateData model.MenuInput) (*model.Menu, error) {
	var newUser = new(model.Menu)
	newUser.Name = updateData.Name
	newUser.Category = updateData.Category
	newUser.Price = updateData.Price
	newUser.Description = updateData.Description

	result, err := service.repository.Update(id, newUser)
	if err != nil {
		return nil, errors.New("error updating")
	}

	return result, nil
}

func (service *menuService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		return errors.New("error deleting")
	}

	return nil
}
