package service

import (
	"errors"
	"mime/multipart"
	"restoran/features/menu/model"
	"restoran/features/menu/repository"
	"restoran/helper"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type MenuServiceInterface interface {
	Insert(fileHeader *multipart.FileHeader, newData model.MenuInput) (*model.MenuResponse, error)
	GetAll(pagination model.QueryParam) ([]model.MenuResponse, *model.Pagination, error)
	GetCategory(queryParam model.QueryParam) ([]model.MenuResponse, *model.Pagination, error)
	GetFavorite() ([]model.Favorite, error)
	GetByName(name string) (*model.MenuResponse, error)
	Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.MenuResponse, error)
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

func (service *menuService) Insert(fileHeader *multipart.FileHeader, newData model.MenuInput) (*model.MenuResponse, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var findMenu = service.repository.GetByName(newData.Name)
	if findMenu != nil {
		return nil, errors.New("menu already exists")
	}

	urlImage, err := service.repository.UploadImage(fileHeader, newData.Name)
	if err != nil {
		return nil, errors.New("upload image failed")
	}

	newData.Image = urlImage

	var newMenu = helper.RequestToMenu(newData)
	result, err := service.repository.Insert(newMenu)
	if err != nil {
		return nil, errors.New("insert data menu failed")
	}

	var menuResponse = helper.MenuToResponse(result)

	return &menuResponse, nil
}

func (service *menuService) GetAll(pagination model.QueryParam) ([]model.MenuResponse, *model.Pagination, error) {
	result, err := service.repository.GetAll(pagination)
	if err != nil {
		return nil, nil, errors.New("get data menu failed")
	}

	var menuResponse []model.MenuResponse

	for _, menu := range result {
		menuResponse = append(menuResponse, helper.MenuToResponse(&menu))
	}

	total, err := service.repository.TotalData()
	if err != nil {
		return nil, nil, errors.New("get total menu failed")
	}

	var paginationResponse = helper.QueryParamToPagination(pagination, total)

	return menuResponse, paginationResponse, nil
}

func (service *menuService) GetCategory(queryParam model.QueryParam) ([]model.MenuResponse, *model.Pagination, error) {
	result, err := service.repository.GetCategory(queryParam)
	if err != nil {
		logrus.Error("Service: Get data by category failed,", err)
		return nil, nil, errors.New("get data menu by category failed")
	}

	var menuResponse []model.MenuResponse
	for _, menu := range result {
		menuResponse = append(menuResponse, helper.MenuToResponse(&menu))
	}

	total, err := service.repository.TotalDataByCategory(queryParam.Category)
	if err != nil {
		return nil, nil, errors.New("get total menu by category failed")
	}

	var paginationResponse = helper.QueryParamToPagination(queryParam, total)

	return menuResponse, paginationResponse, nil
}

func (service *menuService) GetFavorite() ([]model.Favorite, error) {
	result, err := service.repository.GetFavorite()
	if err != nil {
		logrus.Error("Service: Get favorite data failed,", err)
		return nil, errors.New("get data favorite menu failed")
	}

	return result, nil
}

func (service *menuService) GetByName(name string) (*model.MenuResponse, error) {
	result := service.repository.GetByName(name)
	if result == nil {
		return nil, errors.New("menu not found")
	}

	var menuResponse = helper.MenuToResponse(result)

	return &menuResponse, nil
}

func (service *menuService) Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.MenuResponse, error) {
	err := service.validator.Struct(updateData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	urlImage, err := service.repository.UploadImage(fileHeader, updateData.Name)
	if err != nil {
		logrus.Error("Service: Upload image failed,", err)
		return nil, errors.New("cannot upload image")
	}

	updateData.Image = urlImage

	var updateMenu = helper.RequestToMenu(updateData)
	result, err := service.repository.Update(id, updateMenu)
	if err != nil {
		logrus.Error("Service: Update data failed: ", err)
		return nil, errors.New("cannot update data")
	}

	var menuResponse = helper.MenuToResponse(result)

	return &menuResponse, nil
}

func (service *menuService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		logrus.Error("Service: Delete data failed: ", err)
		return err
	}

	return nil
}
