package service

import (
	"errors"
	"restoran/features/transaction/mocks"
	"restoran/features/transaction/model"
	"restoran/helper"
	mockHelper "restoran/helper/mocks"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestInsertWithPaymentGateway(t *testing.T) {
	var repository = mocks.NewTransactionRepositoryInterface(t)
	var validate = validator.New()
	var generator = mockHelper.NewGeneratorInterface(t)
	var service = NewTransactionService(repository, validate, generator)

	var newData = model.TransactionInput{
		OrderID: 1,
	}

	var invalidData = model.TransactionInput{}

	var transactionData = model.Transaction{
		ID:      1,
		OrderID: "randomOrderID",
		Status:  "Pending",
	}

	var orderData = model.Order{
		ID:    1,
		Total: 10000,
	}

	var token = "randomToken"
	var redirectURL = "http://randomredirecturl.com"

	t.Run("Success insert", func(t *testing.T) {
		var newTransaction = helper.RequestToTransaction(newData)
		newTransaction.OrderID = "randomUUID"

		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		repository.On("Insert", newTransaction).Return(&transactionData, nil).Once()
		repository.On("GetOrder", int(transactionData.ID)).Return(&orderData, nil).Once()
		repository.On("SnapRequest", transactionData.OrderID, int64(orderData.Total)).Return(token, redirectURL).Once()

		result, err := service.InsertWithPaymentGateway(newData)
		assert.Nil(t, err)
		assert.Equal(t, newData.OrderID, result.ID)
		generator.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Validation failed", func(t *testing.T) {
		result, err := service.InsertWithPaymentGateway(invalidData)

		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Generate id failed", func(t *testing.T) {
		generator.On("GenerateUUID").Return("", errors.New("id generator failed")).Once()

		result, err := service.InsertWithPaymentGateway(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "order id generator failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
	})

	t.Run("Insert failed", func(t *testing.T) {
		var newTransaction = helper.RequestToTransaction(newData)
		newTransaction.OrderID = "randomUUID"

		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		repository.On("Insert", newTransaction).Return(nil, errors.New("insert data failed")).Once()

		result, err := service.InsertWithPaymentGateway(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "insert data transaction failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Get order failed", func(t *testing.T) {
		var newTransaction = helper.RequestToTransaction(newData)
		newTransaction.OrderID = "randomUUID"

		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		repository.On("Insert", newTransaction).Return(&transactionData, nil).Once()
		repository.On("GetOrder", int(transactionData.ID)).Return(nil, errors.New("get order error")).Once()

		result, err := service.InsertWithPaymentGateway(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "get data order failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Snap request failed", func(t *testing.T) {
		var newTransaction = helper.RequestToTransaction(newData)
		newTransaction.OrderID = "randomUUID"

		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		repository.On("Insert", newTransaction).Return(&transactionData, nil).Once()
		repository.On("GetOrder", int(transactionData.ID)).Return(&orderData, nil).Once()
		repository.On("SnapRequest", transactionData.OrderID, int64(orderData.Total)).Return("", "").Once()

		result, err := service.InsertWithPaymentGateway(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "create transaction failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
		repository.AssertExpectations(t)
	})
}

func TestInsertWithoutPaymentGateway(t *testing.T) {
	var repository = mocks.NewTransactionRepositoryInterface(t)
	var validate = validator.New()
	var generator = mockHelper.NewGeneratorInterface(t)
	var service = NewTransactionService(repository, validate, generator)

	var newData = model.TransactionInput{
		OrderID:       1,
		PaymentMethod: "cash",
	}

	var invalidData = model.TransactionInput{}

	var transactionData = model.Transaction{
		ID:            1,
		OrderID:       "randomOrderID",
		PaymentMethod: "cash",
		Status:        "success",
	}

	t.Run("Success insert data", func(t *testing.T) {
		var newTransaction = helper.RequestToTransaction(newData)
		newTransaction.OrderID = "randomUUID"
		newTransaction.PaymentMethod = "cash"
		newTransaction.Status = "success"

		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		repository.On("Insert", newTransaction).Return(&transactionData, nil).Once()
		repository.On("UpdateStatusOrder", transactionData.ID, "Paid").Return(nil).Once()

		result, err := service.InsertWithoutPaymentGateway(newData)
		assert.Nil(t, err)
		assert.Equal(t, newData.OrderID, result.ID)
		assert.Equal(t, newData.PaymentMethod, result.PaymentMethod)
		generator.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Validation failed", func(t *testing.T) {
		result, err := service.InsertWithoutPaymentGateway(invalidData)

		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Generate id failed", func(t *testing.T) {
		generator.On("GenerateUUID").Return("", errors.New("id generator failed")).Once()

		result, err := service.InsertWithoutPaymentGateway(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "order id generator failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
	})

	t.Run("Insert data failed", func(t *testing.T) {
		var newTransaction = helper.RequestToTransaction(newData)
		newTransaction.OrderID = "randomUUID"
		newTransaction.PaymentMethod = "cash"
		newTransaction.Status = "success"

		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		repository.On("Insert", newTransaction).Return(nil, errors.New("insert data failed")).Once()

		result, err := service.InsertWithoutPaymentGateway(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "insert data transaction failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Update status order error", func(t *testing.T) {
		var newTransaction = helper.RequestToTransaction(newData)
		newTransaction.OrderID = "randomUUID"
		newTransaction.PaymentMethod = "cash"
		newTransaction.Status = "success"

		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		repository.On("Insert", newTransaction).Return(&transactionData, nil).Once()
		repository.On("UpdateStatusOrder", transactionData.ID, "Paid").Return(errors.New("update status order error")).Once()

		result, err := service.InsertWithoutPaymentGateway(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "update status order failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
		repository.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	var repository = mocks.NewTransactionRepositoryInterface(t)
	var validate = validator.New()
	var generator = mockHelper.NewGeneratorInterface(t)
	var service = NewTransactionService(repository, validate, generator)

	var pagination = model.QueryParam{
		Page:     1,
		PageSize: 2,
	}

	var transactionList = []model.Transaction{
		{
			ID:      1,
			OrderID: "randomUser",
			Status:  "Pending",
		},
		{
			ID:      2,
			OrderID: "randomUser",
			Status:  "Pending",
		},
	}

	t.Run("Success get all data", func(t *testing.T) {
		repository.On("GetAll", pagination).Return(transactionList, nil).Once()

		result, err := service.GetAll(pagination)
		assert.Nil(t, err)
		assert.Equal(t, len(transactionList), len(result))
		assert.Equal(t, transactionList[0].ID, result[0].ID)
		repository.AssertExpectations(t)
	})

	t.Run("Get all data failed", func(t *testing.T) {
		repository.On("GetAll", pagination).Return(nil, errors.New("Get all data failed")).Once()

		result, err := service.GetAll(pagination)
		assert.Error(t, err)
		assert.EqualError(t, err, "get data transaction failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	var repository = mocks.NewTransactionRepositoryInterface(t)
	var validate = validator.New()
	var generator = mockHelper.NewGeneratorInterface(t)
	var service = NewTransactionService(repository, validate, generator)

	var transactionData = model.Transaction{
		ID:      1,
		OrderID: "randomOrderID",
		Status:  "Pending",
	}

	t.Run("Success get data by id", func(t *testing.T) {
		repository.On("GetByID", 1).Return(&transactionData, nil).Once()

		result, err := service.GetByID(1)
		assert.Nil(t, err)
		assert.Equal(t, transactionData.ID, result.ID)
		assert.Equal(t, transactionData.OrderID, result.OrderID)
		repository.AssertExpectations(t)
	})

	t.Run("Get data by id failed", func(t *testing.T) {
		repository.On("GetByID", 1).Return(nil, errors.New("get data by id error")).Once()

		result, err := service.GetByID(1)
		assert.Error(t, err)
		assert.EqualError(t, err, "get data transaction by id failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	var repository = mocks.NewTransactionRepositoryInterface(t)
	var validate = validator.New()
	var generator = mockHelper.NewGeneratorInterface(t)
	var service = NewTransactionService(repository, validate, generator)

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

func TestNotifications(t *testing.T) {
	var repository = mocks.NewTransactionRepositoryInterface(t)
	var validate = validator.New()
	var generator = mockHelper.NewGeneratorInterface(t)
	var service = NewTransactionService(repository, validate, generator)

	var notificationsPayload = map[string]any{
		"order_id":     "randomOrderID",
		"payment_type": "qris",
	}
	var invalidNotificationsPayload = map[string]any{}

	var paymentMethodNotFound = map[string]any{
		"order_id": "randomOrderID",
	}

	var status = model.Status{
		Status:        "Success",
		Order:         "Paid",
		PaymentMethod: "qris",
	}

	var transactionData = model.Transaction{
		ID:      1,
		OrderID: "randomOrderID",
	}

	t.Run("Success notifications", func(t *testing.T) {
		repository.On("CheckTransaction", notificationsPayload["order_id"]).Return(status, nil).Once()
		repository.On("GetByOrderID", notificationsPayload["order_id"]).Return(&transactionData, nil).Once()
		repository.On("UpdateStatusTransaction", transactionData.ID, status).Return(nil).Once()
		repository.On("UpdateStatusOrder", transactionData.ID, status.Order).Return(nil).Once()

		err := service.Notifications(notificationsPayload)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Order id not found", func(t *testing.T) {
		err := service.Notifications(invalidNotificationsPayload)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid notification payload")
	})

	t.Run("Payment method not found", func(t *testing.T) {
		err := service.Notifications(paymentMethodNotFound)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid notification payload")
	})

	t.Run("Check transaction failed", func(t *testing.T) {
		repository.On("CheckTransaction", notificationsPayload["order_id"]).Return(model.Status{}, errors.New("check transaction error")).Once()

		err := service.Notifications(notificationsPayload)
		assert.Error(t, err)
		assert.EqualError(t, err, "check transaction error")
		repository.AssertExpectations(t)
	})

	t.Run("Get by order id failed", func(t *testing.T) {
		repository.On("CheckTransaction", notificationsPayload["order_id"]).Return(status, nil).Once()
		repository.On("GetByOrderID", notificationsPayload["order_id"]).Return(nil, errors.New("get by order id error")).Once()

		err := service.Notifications(notificationsPayload)
		assert.Error(t, err)
		assert.EqualError(t, err, "transaction data not found")
		repository.AssertExpectations(t)
	})

	t.Run("Update status transaction failed", func(t *testing.T) {
		repository.On("CheckTransaction", notificationsPayload["order_id"]).Return(status, nil).Once()
		repository.On("GetByOrderID", notificationsPayload["order_id"]).Return(&transactionData, nil).Once()
		repository.On("UpdateStatusTransaction", transactionData.ID, status).Return(errors.New("update status transaction error")).Once()

		err := service.Notifications(notificationsPayload)
		assert.Error(t, err)
		assert.EqualError(t, err, "update status transaction error")
		repository.AssertExpectations(t)
	})

	t.Run("Update status order error", func(t *testing.T) {
		repository.On("CheckTransaction", notificationsPayload["order_id"]).Return(status, nil).Once()
		repository.On("GetByOrderID", notificationsPayload["order_id"]).Return(&transactionData, nil).Once()
		repository.On("UpdateStatusTransaction", transactionData.ID, status).Return(nil).Once()
		repository.On("UpdateStatusOrder", transactionData.ID, status.Order).Return(errors.New("update status order error")).Once()

		err := service.Notifications(notificationsPayload)
		assert.Error(t, err)
		assert.EqualError(t, err, "update status order error")
		repository.AssertExpectations(t)
	})
}
