package service

import (
	"errors"
	"restoran/features/order/model"
	"restoran/features/order/repository"
	"restoran/helper"

	"github.com/sirupsen/logrus"
)

type OrderServiceInterface interface {
	Insert(newData model.OrderInput) (*model.Order, error)
	GetAll() ([]model.Order, error)
	GetByID(id int) (*model.Order, error)
	Update(id int, updateData model.OrderInput) (*model.Order, error)
	Delete(id int) error
}

type orderService struct {
	repository repository.OrderRepositoryInterface
}

func NewOrderService(repo repository.OrderRepositoryInterface) OrderServiceInterface {
	return &orderService{
		repository: repo,
	}
}

func (service *orderService) Insert(newData model.OrderInput) (*model.Order, error) {
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
		subTotal := price[i] * qty.Quantity
		totalPrice += subTotal
	}

	var newOrder = helper.RequestToOrder(newData, totalPrice)

	result, err := service.repository.Insert(newOrder)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("Cannot insert data " + err.Error())
	}

	return result, nil
}

func (service *orderService) GetAll() ([]model.Order, error) {
	result, err := service.repository.GetAll()
	if err != nil {
		logrus.Error("Service: Get all data failed,", err)
		return nil, errors.New("Cannot get all data " + err.Error())
	}

	return result, nil
}

func (service *orderService) GetByID(id int) (*model.Order, error) {
	result, err := service.repository.GetByID(id)
	if err != nil {
		logrus.Error("Service: Get data by id failed,", err)
		return nil, errors.New("Cannot get data by id " + err.Error())
	}

	return result, nil
}

func (service *orderService) Update(id int, updateData model.OrderInput) (*model.Order, error) {
	var updateUser = helper.RequestToOrder(updateData, 0)
	result, err := service.repository.Update(id, updateUser)
	if err != nil {
		logrus.Error("Service: Update data failed: ", err)
		return nil, errors.New("Cannot update data " + err.Error())
	}

	return result, nil
}

func (service *orderService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		logrus.Error("Service: Delete data failed: ", err)
		return errors.New("Cannot delete data " + err.Error())
	}

	return nil
}
