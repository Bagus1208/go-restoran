package service

import (
	"errors"
	"restoran/features/order/model"
	"restoran/features/order/repository"
	"restoran/helper"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type OrderServiceInterface interface {
	Insert(newData model.OrderInput) (*model.Order, error)
	GetAll(pagination model.Pagination) ([]model.Order, error)
	GetByID(id int) (*model.Order, error)
	Delete(id int) error
}

type orderService struct {
	repository repository.OrderRepositoryInterface
	validator  *validator.Validate
}

func NewOrderService(repo repository.OrderRepositoryInterface, validate *validator.Validate) OrderServiceInterface {
	return &orderService{
		repository: repo,
		validator:  validate,
	}
}

func (service *orderService) Insert(newData model.OrderInput) (*model.Order, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var menuName []string
	for _, order := range newData.Orders {
		menuName = append(menuName, order.MenuName)
	}

	findMenu, price := service.repository.FindMenu(menuName)
	if !findMenu {
		return nil, errors.New("menu not found")
	}

	var totalPrice int
	for i, qty := range newData.Orders {
		totalPrice += price[i] * qty.Quantity
	}

	var newOrder = helper.RequestToOrder(newData, totalPrice)
	result, err := service.repository.Insert(newOrder)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("cannot insert data: " + err.Error())
	}

	return result, nil
}

func (service *orderService) GetAll(pagination model.Pagination) ([]model.Order, error) {
	if pagination.Page <= 0 || pagination.PageSize <= 0 {
		return nil, errors.New("invalid page and page_size value")
	}

	result, err := service.repository.GetAll(pagination)
	if err != nil {
		logrus.Error("Service: Get all data failed,", err)
		return nil, errors.New("cannot get all data: " + err.Error())
	}

	return result, nil
}

func (service *orderService) GetByID(id int) (*model.Order, error) {
	result, err := service.repository.GetByID(id)
	if err != nil {
		logrus.Error("Service: Get data by id failed,", err)
		return nil, errors.New("cannot get data by id: " + err.Error())
	}

	return result, nil
}

func (service *orderService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		logrus.Error("Service: Delete data failed: ", err)
		return errors.New("cannot delete data: " + err.Error())
	}

	return nil
}
