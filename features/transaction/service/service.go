package service

import (
	"errors"
	"restoran/features/transaction/model"
	"restoran/features/transaction/repository"
	"restoran/helper"

	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type TransactionServiceInterface interface {
	Insert(newData model.TransactionInput) (*model.TransactionInputResponse, error)
	GetAll(pagination model.QueryParam) ([]model.TransactionResponse, error)
	GetByID(id int) (*model.TransactionResponse, error)
	Delete(id int) error
	Notifications(notificationPayload map[string]any) error
}

type transactionService struct {
	repository    repository.TransactionRepositoryInterface
	validator     *validator.Validate
	snapClient    snap.Client
	coreAPIClient coreapi.Client
}

func NewTransactionService(repo repository.TransactionRepositoryInterface, validate *validator.Validate, snapclient snap.Client, coreAPIClient coreapi.Client) TransactionServiceInterface {
	return &transactionService{
		repository:    repo,
		validator:     validate,
		snapClient:    snapclient,
		coreAPIClient: coreAPIClient,
	}
}

func (service *transactionService) Insert(newData model.TransactionInput) (*model.TransactionInputResponse, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var newTransaction = helper.RequestToTransaction(newData)
	newTransaction.OrderID = helper.GenerateUUID()

	result, err := service.repository.Insert(newTransaction)
	if err != nil {
		return nil, errors.New("insert data transaction failed")
	}

	order, err := service.repository.GetOrder(int(result.ID))
	if err != nil {
		return nil, errors.New("get data order failed")
	}

	snapResponse, _ := helper.CreateSnapRequest(service.snapClient, result.OrderID, int64(order.Total))
	var transactionInputResponse = helper.TransactionToResponseInput(result, snapResponse.Token, snapResponse.RedirectURL)

	return transactionInputResponse, nil
}

func (service *transactionService) GetAll(pagination model.QueryParam) ([]model.TransactionResponse, error) {
	if pagination.Page <= 0 || pagination.PageSize <= 0 {
		return nil, errors.New("invalid page and page_size value")
	}

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
	transactionID, exist := notificationPayload["order_id"].(string)
	if !exist {
		return errors.New("invalid notification payload")
	}

	transactionStatusResp, err := service.coreAPIClient.CheckTransaction(transactionID)
	if err != nil {
		return err
	} else {
		if transactionStatusResp != nil {
			var status = helper.TransactionStatus(transactionStatusResp)

			transaction, _ := service.repository.GetByOrderID(transactionID)

			err := service.repository.UpdateStatusTransaction(int(transaction.ID), status.Transaction)
			if err != nil {
				return err
			}

			err = service.repository.UpdateStatusOrder(int(transaction.ID), status.Order)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
