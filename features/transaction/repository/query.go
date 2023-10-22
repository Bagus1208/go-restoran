package repository

import (
	"errors"
	"restoran/features/transaction/model"

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
	UpdateStatusTransaction(id int, status string) error
	UpdateStatusOrder(id int, status string) error
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(DB *gorm.DB) TransactionRepositoryInterface {
	return &transactionRepo{
		db: DB,
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

func (repository *transactionRepo) UpdateStatusTransaction(id int, status string) error {
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

func (repository *transactionRepo) UpdateStatusOrder(id int, status string) error {
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
