package service

import (
	"context"
	"errors"
	"mime/multipart"
	"restoran/features/menu/model"
	"restoran/features/menu/repository"
	"restoran/helper"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type MenuServiceInterface interface {
	Insert(fileHeader *multipart.FileHeader, newData model.MenuInput) (*model.Menu, error)
	GetAll(model.Pagination) ([]model.Menu, error)
	GetCategory(category string, pagination model.Pagination) ([]model.Menu, error)
	Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.Menu, error)
	Delete(id int) error
}

type menuService struct {
	repository repository.MenuRepositoryInterface
	validator  *validator.Validate
}

func NewMenuService(repo repository.MenuRepositoryInterface, validate *validator.Validate) MenuServiceInterface {
	return &menuService{
		repository: repo,
		validator:  validate,
	}
}

func (service *menuService) Insert(fileHeader *multipart.FileHeader, newData model.MenuInput) (*model.Menu, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var findMenu = service.repository.GetByName(newData.Name)
	if findMenu != nil {
		return nil, errors.New("menu already exists")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := fileHeader.Open()
	if err != nil {
		logrus.Error("Service: Open fileHeader failed,", err)
		return nil, errors.New("cannot open fileHeader " + err.Error())
	}

	urlImage, err := service.repository.UploadImage(ctx, file, newData.Name)
	if err != nil {
		logrus.Error("Service: Upload image failed,", err)
		return nil, errors.New("cannot upload image " + err.Error())
	}

	newData.Image = urlImage

	var newMenu = helper.RequestToMenu(newData)
	result, err := service.repository.Insert(newMenu)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("cannot insert data " + err.Error())
	}

	return result, nil
}

func (service *menuService) GetAll(pagination model.Pagination) ([]model.Menu, error) {
	if pagination.Page <= 0 || pagination.PageSize <= 0 {
		return nil, errors.New("invalid page and pageSize value")
	}

	result, err := service.repository.GetAll(pagination)
	if err != nil {
		logrus.Error("Service: Get all data failed,", err)
		return nil, errors.New("cannot get all data " + err.Error())
	}

	return result, nil
}

func (service *menuService) GetCategory(category string, pagination model.Pagination) ([]model.Menu, error) {
	if pagination.Page <= 0 || pagination.PageSize <= 0 {
		return nil, errors.New("invalid page and pageSize value")
	}

	result, err := service.repository.GetCategory(category, pagination)
	if err != nil {
		logrus.Error("Service: Get data by category failed,", err)
		return nil, errors.New("cannot get data by category " + err.Error())
	}

	return result, nil
}

func (service *menuService) Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.Menu, error) {
	err := service.validator.Struct(updateData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var findMenu = service.repository.GetByName(updateData.Name)
	if findMenu != nil {
		return nil, errors.New("menu already exists")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := fileHeader.Open()
	if err != nil {
		logrus.Error("Service: Open fileHeader failed,", err)
		return nil, errors.New("cannot open fileHeader " + err.Error())
	}

	urlImage, err := service.repository.UploadImage(ctx, file, updateData.Name)
	if err != nil {
		logrus.Error("Service: Upload image failed,", err)
		return nil, errors.New("cannot upload image " + err.Error())
	}

	updateData.Image = urlImage

	var updateMenu = helper.RequestToMenu(updateData)
	result, err := service.repository.Update(id, updateMenu)
	if err != nil {
		logrus.Error("Service: Update data failed: ", err)
		return nil, errors.New("cannot update data " + err.Error())
	}

	return result, nil
}

func (service *menuService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		logrus.Error("Service: Delete data failed: ", err)
		return errors.New("cannot delete data " + err.Error())
	}

	return nil
}
