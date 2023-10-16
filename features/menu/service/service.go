package service

import (
	"context"
	"errors"
	"mime/multipart"
	"restoran/features/menu/model"
	"restoran/features/menu/repository"
	"restoran/helper"
	"time"

	"github.com/sirupsen/logrus"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := fileHeader.Open()
	if err != nil {
		logrus.Error("Service: Open fileHeader failed,", err)
		return nil, errors.New("Cannot open fileHeader " + err.Error())
	}

	urlImage, err := service.repository.UploadImage(ctx, file, newData.Name)
	if err != nil {
		logrus.Error("Service: Upload image failed,", err)
		return nil, errors.New("Cannot upload image " + err.Error())
	}

	newData.Image = urlImage

	var newMenu = helper.RequestToMenu(newData)
	result, err := service.repository.Insert(newMenu)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("Cannot insert data " + err.Error())
	}

	return result, nil
}

func (service *menuService) GetAll() ([]model.Menu, error) {
	result, err := service.repository.GetAll()
	if err != nil {
		logrus.Error("Service: Get all data failed,", err)
		return nil, errors.New("Cannot get all data " + err.Error())
	}

	return result, nil
}

func (service *menuService) GetCategory(category string) ([]model.Menu, error) {
	result, err := service.repository.GetCategory(category)
	if err != nil {
		logrus.Error("Service: Get data by category failed,", err)
		return nil, errors.New("Cannot get data by category " + err.Error())
	}

	return result, nil
}

func (service *menuService) Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.Menu, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := fileHeader.Open()
	if err != nil {
		logrus.Error("Service: Open fileHeader failed,", err)
		return nil, errors.New("Cannot open fileHeader " + err.Error())
	}

	urlImage, err := service.repository.UploadImage(ctx, file, updateData.Name)
	if err != nil {
		logrus.Error("Service: Upload image failed,", err)
		return nil, errors.New("Cannot upload image " + err.Error())
	}

	updateData.Image = urlImage

	var updateMenu = helper.RequestToMenu(updateData)
	result, err := service.repository.Update(id, updateMenu)
	if err != nil {
		logrus.Error("Service: Update data failed: ", err)
		return nil, errors.New("Cannot update data " + err.Error())
	}

	return result, nil
}

func (service *menuService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		logrus.Error("Service: Delete data failed: ", err)
		return errors.New("Cannot delete data " + err.Error())
	}

	return nil
}
