package service

import (
	"errors"
	"restoran/features/order/mocks"
	"restoran/features/order/model"
	"restoran/helper"
	mockHelper "restoran/helper/mocks"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	var jwt = mockHelper.NewJWTInterface(t)
	var validate = validator.New()
	var repository = mocks.NewOrderRepositoryInterface(t)
	var service = NewOrderService(repository, validate, jwt)

	var newData = model.OrderInput{
		Orders: []model.OrderDetail{
			{
				MenuID:   1,
				Quantity: 5,
			},
		},
	}

	var invalidData = model.OrderInput{}

	var order = model.Order{
		ID:      1,
		NoTable: 1,
		Orders: []model.OrderDetail{
			{
				ID:       1,
				MenuID:   1,
				Quantity: 5,
			},
		},
		Total:  10000,
		Status: "unpaid",
	}

	var tokenString = "randomtokenString"
	var menuID = []int{1}

	t.Run("Success insert data", func(t *testing.T) {
		var newOrder = helper.RequestToOrder(newData, 10000)
		newOrder.NoTable = 1

		jwt.On("ExtractToken", tokenString).Return(order.NoTable, nil).Once()
		repository.On("FindMenu", menuID).Return(true, []int{2000}).Once()
		repository.On("Insert", newOrder).Return(&order, nil).Once()

		result, err := service.Insert(newData, tokenString)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newData.Orders[0].MenuID, result.Orders[0].MenuID)
		assert.Equal(t, newData.Orders[0].Quantity, result.Orders[0].Quantity)
		repository.AssertExpectations(t)
	})

	t.Run("Get no table failed", func(t *testing.T) {
		jwt.On("ExtractToken", tokenString).Return(0, errors.New("token error")).Once()

		result, err := service.Insert(newData, tokenString)
		assert.Error(t, err)
		assert.EqualError(t, err, "token error")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Validation failed", func(t *testing.T) {
		jwt.On("ExtractToken", tokenString).Return(order.NoTable, nil).Once()

		result, err := service.Insert(invalidData, tokenString)

		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Find menu failed", func(t *testing.T) {
		jwt.On("ExtractToken", tokenString).Return(order.NoTable, nil).Once()
		repository.On("FindMenu", menuID).Return(false, nil).Once()

		result, err := service.Insert(newData, tokenString)
		assert.Error(t, err)
		assert.EqualError(t, err, "menu not found")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Insert data failed", func(t *testing.T) {
		var newOrder = helper.RequestToOrder(newData, 10000)
		newOrder.NoTable = 1

		jwt.On("ExtractToken", tokenString).Return(order.NoTable, nil).Once()
		repository.On("FindMenu", menuID).Return(true, []int{2000}).Once()
		repository.On("Insert", newOrder).Return(nil, errors.New("insert data failed")).Once()

		result, err := service.Insert(newData, tokenString)
		assert.Error(t, err)
		assert.EqualError(t, err, "insert data order failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	var jwt = mockHelper.NewJWTInterface(t)
	var validate = validator.New()
	var repository = mocks.NewOrderRepositoryInterface(t)
	var service = NewOrderService(repository, validate, jwt)

	var pagination = model.Pagination{
		Page:     1,
		PageSize: 2,
	}

	var orderList = []model.Order{
		{
			ID:      1,
			NoTable: 1,
			Orders: []model.OrderDetail{
				{
					ID:       1,
					MenuID:   1,
					Quantity: 5,
				},
			},
			Total:  10000,
			Status: "unpaid",
		},
		{
			ID:      2,
			NoTable: 2,
			Orders: []model.OrderDetail{
				{
					ID:       2,
					MenuID:   1,
					Quantity: 5,
				},
			},
			Total:  10000,
			Status: "unpaid",
		},
	}

	t.Run("Success get all data", func(t *testing.T) {
		repository.On("GetAll", pagination).Return(orderList, nil).Once()

		result, err := service.GetAll(pagination)
		assert.Nil(t, err)
		assert.Equal(t, len(orderList), len(result))
		assert.Equal(t, orderList[0].ID, result[0].ID)
		repository.AssertExpectations(t)
	})

	t.Run("Get all data failed", func(t *testing.T) {
		repository.On("GetAll", pagination).Return(nil, errors.New("get data error")).Once()

		result, err := service.GetAll(pagination)
		assert.Error(t, err)
		assert.EqualError(t, err, "get data order failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	var jwt = mockHelper.NewJWTInterface(t)
	var validate = validator.New()
	var repository = mocks.NewOrderRepositoryInterface(t)
	var service = NewOrderService(repository, validate, jwt)

	var order = model.Order{
		ID:      1,
		NoTable: 1,
		Orders: []model.OrderDetail{
			{
				ID:       1,
				MenuID:   1,
				Quantity: 5,
			},
		},
		Total:  10000,
		Status: "unpaid",
	}

	t.Run("Success get data by id", func(t *testing.T) {
		repository.On("GetByID", 1).Return(&order, nil).Once()

		result, err := service.GetByID(1)
		assert.Nil(t, err)
		assert.Equal(t, order.ID, result.ID)
		repository.AssertExpectations(t)
	})

	t.Run("Get data by id failed", func(t *testing.T) {
		repository.On("GetByID", 1).Return(nil, errors.New("get data by id error")).Once()

		result, err := service.GetByID(1)
		assert.Error(t, err)
		assert.EqualError(t, err, "get data order by id failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	var jwt = mockHelper.NewJWTInterface(t)
	var validate = validator.New()
	var repository = mocks.NewOrderRepositoryInterface(t)
	var service = NewOrderService(repository, validate, jwt)

	t.Run("success delete data", func(t *testing.T) {
		repository.On("Delete", 1).Return(nil).Once()

		err := service.Delete(1)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("delete data failed", func(t *testing.T) {
		repository.On("Delete", 1).Return(errors.New("delete data error")).Once()

		err := service.Delete(1)
		assert.Error(t, err)
		assert.EqualError(t, err, "delete data error")
		repository.AssertExpectations(t)
	})
}
