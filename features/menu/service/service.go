package service

import (
	"errors"
	"mime/multipart"
	"restoran/features/menu/model"
	"restoran/features/menu/repository"
	"restoran/helper"
)

type MenuServiceInterface interface {
	Insert(fileHeader *multipart.FileHeader, newData model.MenuInput) (*model.Menu, error)
	GetAll() ([]model.Menu, error)
	GetCategory(category string) ([]model.Menu, error)
	Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.Menu, error)
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

func (service *menuService) Insert(fileHeader *multipart.FileHeader, newData model.MenuInput) (*model.Menu, error) {
	urlImage, err := helper.UploadImageToCDN(fileHeader, newData.Name)
	if err != nil {
		return nil, errors.New("unprocessable image")
	}
	newData.Image = urlImage

	var newUser = helper.ToMenu(newData)

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

func (service *menuService) Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.Menu, error) {
	urlImage, err := helper.UploadImageToCDN(fileHeader, updateData.Name)
	if err != nil {
		return nil, errors.New("unprocessable image")
	}
	updateData.Image = urlImage

	var newUser = helper.ToMenu(updateData)

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
