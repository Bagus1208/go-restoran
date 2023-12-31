package repository

import (
	"errors"
	"restoran/features/order/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	Insert(newData *model.Order) (*model.Order, error)
	GetAll(pagination model.Pagination) ([]model.Order, error)
	GetByID(id int) (*model.Order, error)
	Delete(id int) error
	FindMenu(menuID []int) (bool, []int)
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(DB *gorm.DB) OrderRepositoryInterface {
	return &orderRepo{
		db: DB,
	}
}

func (repository *orderRepo) Insert(newData *model.Order) (*model.Order, error) {
	result := repository.db.Create(newData)
	if result.Error != nil {
		logrus.Error("Repository: Inserting order error,", result.Error)
		return nil, result.Error
	}

	return newData, nil
}

func (repository *orderRepo) GetAll(pagination model.Pagination) ([]model.Order, error) {
	var orders []model.Order
	var offset = (pagination.Page - 1) * pagination.PageSize

	result := repository.db.Offset(offset).Limit(pagination.PageSize).Preload("Orders").Find(&orders)
	if result.Error != nil {
		logrus.Error("Repository: Get all order error,", result.Error)
		return nil, result.Error
	}

	return orders, nil
}

func (repository *orderRepo) GetByID(id int) (*model.Order, error) {
	var order model.Order

	result := repository.db.Preload("Orders").Where("id =?", id).First(&order)

	if result.Error != nil {
		logrus.Error("Repository: Get order by id error,", result.Error)
		return nil, result.Error
	}

	return &order, nil
}

func (repository *orderRepo) Delete(id int) error {
	var deleteOrder model.Order
	deleteOrder.ID = uint(id)
	result := repository.db.Delete(&deleteOrder)
	if result.Error != nil {
		logrus.Error("Repository: Delete order error,", result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		logrus.Error("Repository: Delete order error,", result.Error)
		return errors.New("data not found")
	}

	return nil
}

func (repository *orderRepo) FindMenu(menuID []int) (bool, []int) {
	var price []int

	result := repository.db.Select("price").Table("menus").Where("id IN ?", menuID).Pluck("price", &price)
	if result.Error != nil {
		logrus.Error("Repository: Find menu error,", result.Error)
		return false, nil
	}

	if len(price) == len(menuID) {
		return true, price
	}

	return false, nil
}
