package service

import (
	"errors"
	"restoran/features/order/model"
	"restoran/features/order/repository"
	"restoran/helper"

	"github.com/go-playground/validator/v10"
)

type OrderServiceInterface interface {
	Insert(newData model.OrderInput, stringToken string) (*model.OrderResponse, error)
	GetAll(pagination model.Pagination) ([]model.OrderResponse, error)
	GetByID(id int) (*model.OrderResponse, error)
	Delete(id int) error
}

type orderService struct {
	repository repository.OrderRepositoryInterface
	validator  *validator.Validate
	jwt        helper.JWTInterface
}

func NewOrderService(repo repository.OrderRepositoryInterface, validate *validator.Validate, jwt helper.JWTInterface) OrderServiceInterface {
	return &orderService{
		repository: repo,
		validator:  validate,
		jwt:        jwt,
	}
}

func (service *orderService) Insert(newData model.OrderInput, stringToken string) (*model.OrderResponse, error) {
	noTable, err := service.jwt.ExtractToken(stringToken)
	if err != nil {
		return nil, err
	}

	newData.NoTable = noTable
	err = service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var menuID []int
	for _, order := range newData.Orders {
		menuID = append(menuID, order.MenuID)
	}

	findMenu, price := service.repository.FindMenu(menuID)
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
		return nil, errors.New("insert data order failed")
	}

	var orderDetailResponse []model.OrderDetailResponse
	for _, order := range result.Orders {
		orderDetailResponse = append(orderDetailResponse, helper.OrderDetailToResponse(&order))
	}

	var orderResponse = helper.OrderToResponse(result)
	orderResponse.Orders = orderDetailResponse

	return &orderResponse, nil
}

func (service *orderService) GetAll(pagination model.Pagination) ([]model.OrderResponse, error) {
	if pagination.Page <= 0 || pagination.PageSize <= 0 {
		return nil, errors.New("invalid page and page_size value")
	}

	result, err := service.repository.GetAll(pagination)
	if err != nil {
		return nil, errors.New("get data order failed")
	}

	var orderResponseList []model.OrderResponse

	for _, order := range result {
		var orderDetailResponse []model.OrderDetailResponse

		for _, orderDetail := range order.Orders {
			orderDetailResponse = append(orderDetailResponse, helper.OrderDetailToResponse(&orderDetail))
		}

		var orderResponse = helper.OrderToResponse(&order)
		orderResponse.Orders = orderDetailResponse
		orderResponseList = append(orderResponseList, orderResponse)
	}

	return orderResponseList, nil
}

func (service *orderService) GetByID(id int) (*model.OrderResponse, error) {
	result, err := service.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("get data order by id failed")
	}

	var orderDetailResponse []model.OrderDetailResponse
	for _, order := range result.Orders {
		orderDetailResponse = append(orderDetailResponse, helper.OrderDetailToResponse(&order))
	}

	var orderResponse = helper.OrderToResponse(result)
	orderResponse.Orders = orderDetailResponse

	return &orderResponse, nil
}

func (service *orderService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
