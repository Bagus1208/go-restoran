package repository

import (
	"errors"
	"restoran/features/transaction/model"
	"restoran/helper"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	Insert(newData *model.Transaction) (*model.Transaction, error)
	GetAll(pagination model.QueryParam) ([]model.Transaction, error)
	GetByID(id int) (*model.Transaction, error)
	GetByOrderID(orderID string) (*model.Transaction, error)
	Delete(id int) error
	GetOrder(id int) (*model.Order, error)
	SnapRequest(orderID string, total int64) (string, string)
	CheckTransaction(orderID string) (model.Status, error)
	UpdateStatusTransaction(id uint, status string) error
	UpdateStatusOrder(id uint, status string) error
}

type transactionRepo struct {
	db            *gorm.DB
	snapClient    snap.Client
	coreAPIClient coreapi.Client
}

func NewTransactionRepo(DB *gorm.DB, snapClient snap.Client, coreAPIClient coreapi.Client) TransactionRepositoryInterface {
	return &transactionRepo{
		db:            DB,
		snapClient:    snapClient,
		coreAPIClient: coreAPIClient,
	}
}

func (repository *transactionRepo) Insert(newData *model.Transaction) (*model.Transaction, error) {
	result := repository.db.Create(newData)
	if result.Error != nil {
		logrus.Error("Repository: Inserting transaction error,", result.Error)
		return nil, result.Error
	}

	return newData, nil
}

func (repository *transactionRepo) GetAll(pagination model.QueryParam) ([]model.Transaction, error) {
	var transactions []model.Transaction
	var offset = (pagination.Page - 1) * pagination.PageSize

	result := repository.db.Offset(offset).Limit(pagination.PageSize).Find(&transactions)
	if result.Error != nil {
		logrus.Error("Repository: Get all transaction error,", result.Error)
		return nil, result.Error
	}

	return transactions, nil
}

func (repository *transactionRepo) GetByID(id int) (*model.Transaction, error) {
	var transaction model.Transaction

	result := repository.db.Where("id = ?", id).First(&transaction)

	if result.Error != nil {
		logrus.Error("Repository: Get transaction by id error,", result.Error)
		return nil, result.Error
	}

	return &transaction, nil
}

func (repository *transactionRepo) GetByOrderID(orderID string) (*model.Transaction, error) {
	var transaction model.Transaction

	result := repository.db.Where("order_id = ?", orderID).First(&transaction)

	if result.Error != nil {
		logrus.Error("Repository: Get transaction by id error,", result.Error)
		return nil, result.Error
	}

	return &transaction, nil
}

func (repository *transactionRepo) Delete(id int) error {
	var deleteTransaction model.Transaction
	deleteTransaction.ID = uint(id)
	result := repository.db.Delete(&deleteTransaction)
	if result.Error != nil {
		logrus.Error("Repository: Delete transaction error,", result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		logrus.Error("Repository: Delete transaction error,", result.Error)
		return errors.New("data not found")
	}

	return nil
}

func (repository *transactionRepo) GetOrder(id int) (*model.Order, error) {
	var order model.Order

	result := repository.db.Table("orders").Select("id, total").Where("id = ?", id).First(&order)

	if result.Error != nil {
		logrus.Error("Repository: Get order by id error,", result.Error)
		return nil, result.Error
	}

	return &order, nil
}

func (repository *transactionRepo) SnapRequest(orderID string, total int64) (string, string) {
	snapResponse, err := helper.CreateSnapRequest(repository.snapClient, orderID, total)
	if err != nil {
		return "", ""
	}

	return snapResponse.Token, snapResponse.RedirectURL
}

func (repository *transactionRepo) CheckTransaction(orderID string) (model.Status, error) {
	var status model.Status

	transactionStatusResp, err := repository.coreAPIClient.CheckTransaction(orderID)
	if err != nil {
		return model.Status{}, err
	} else {
		if transactionStatusResp != nil {
			status = helper.TransactionStatus(transactionStatusResp)
			return status, nil
		}
	}

	return model.Status{}, err
}

func (repository *transactionRepo) UpdateStatusTransaction(id uint, status string) error {
	result := repository.db.Table("transactions").Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		logrus.Error("Repository: Update transaction status error,", result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		logrus.Error("Repository: No row Affected ,", result.Error)
		return errors.New("data not found")
	}

	return nil
}

func (repository *transactionRepo) UpdateStatusOrder(id uint, status string) error {
	result := repository.db.Table("orders").Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		logrus.Error("Repository: Update order status error,", result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		logrus.Error("Repository: Update order status error,", result.Error)
		return errors.New("data not found")
	}

	return nil
}
