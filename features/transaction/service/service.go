package service

import (
	"errors"
	"restoran/features/transaction/model"
	"restoran/features/transaction/repository"
	"restoran/helper"

	"github.com/go-playground/validator/v10"
)

type TransactionServiceInterface interface {
	Insert(newData model.TransactionInput) (*model.TransactionInputResponse, error)
	GetAll(pagination model.QueryParam) ([]model.TransactionResponse, error)
	GetByID(id int) (*model.TransactionResponse, error)
	Delete(id int) error
	Notifications(notificationPayload map[string]any) error
}

type transactionService struct {
	repository repository.TransactionRepositoryInterface
	validator  *validator.Validate
	generate   helper.GeneratorInterface
}

func NewTransactionService(repo repository.TransactionRepositoryInterface,
	validate *validator.Validate,
	generator helper.GeneratorInterface) TransactionServiceInterface {
	return &transactionService{
		repository: repo,
		validator:  validate,
		generate:   generator,
	}
}

func (service *transactionService) Insert(newData model.TransactionInput) (*model.TransactionInputResponse, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var newTransaction = helper.RequestToTransaction(newData)
	newTransaction.OrderID, err = service.generate.GenerateUUID()
	if err != nil {
		return nil, errors.New("order id generator failed")
	}

	result, err := service.repository.Insert(newTransaction)
	if err != nil {
		return nil, errors.New("insert data transaction failed")
	}

	order, err := service.repository.GetOrder(int(result.ID))
	if err != nil {
		return nil, errors.New("get data order failed")
	}

	token, redirectURL := service.repository.SnapRequest(result.OrderID, int64(order.Total))
	if token == "" && redirectURL == "" {
		return nil, errors.New("create transaction failed")
	}

	var transactionInputResponse = helper.TransactionToResponseInput(result, token, redirectURL)

	return transactionInputResponse, nil
}

func (service *transactionService) GetAll(pagination model.QueryParam) ([]model.TransactionResponse, error) {
	result, err := service.repository.GetAll(pagination)
	if err != nil {
		return nil, errors.New("get data transaction failed")
	}

	var transactionResponse []model.TransactionResponse

	for _, transaction := range result {
		transactionResponse = append(transactionResponse, *helper.TransactionToResponse(&transaction))
	}

	return transactionResponse, nil
}

func (service *transactionService) GetByID(id int) (*model.TransactionResponse, error) {
	result, err := service.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("get data transaction by id failed")
	}

	var transactionResponse = helper.TransactionToResponse(result)

	return transactionResponse, nil
}

func (service *transactionService) Delete(id int) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (service *transactionService) Notifications(notificationPayload map[string]any) error {
	orderID, exist := notificationPayload["order_id"].(string)
	if !exist {
		return errors.New("invalid notification payload")
	}

	status, err := service.repository.CheckTransaction(orderID)
	if err != nil {
		return err
	}

	transaction, err := service.repository.GetByOrderID(orderID)
	if err != nil {
		return errors.New("transaction data not found")
	}

	err = service.repository.UpdateStatusTransaction(transaction.ID, status.Transaction)
	if err != nil {
		return err
	}

	err = service.repository.UpdateStatusOrder(transaction.ID, status.Order)
	if err != nil {
		return err
	}

	return nil
}
