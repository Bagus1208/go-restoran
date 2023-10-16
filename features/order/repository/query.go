package repository

import (
	"restoran/features/order/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	Insert(newData *model.Order) (*model.Order, error)
	GetAll() ([]model.Order, error)
	GetByID(id int) (*model.Order, error)
	Delete(id int) error
	FindMenu(menuNames []string) (bool, []int)
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

func (repository *orderRepo) GetAll() ([]model.Order, error) {
	var orders []model.Order

	result := repository.db.Preload("Orders").Find(&orders)

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

	return nil
}

func (repository *orderRepo) FindMenu(menuNames []string) (bool, []int) {
	var price []int

	result := repository.db.Select("price").Table("menus").Where("name IN ?", menuNames).Pluck("price", &price)
	if result.Error != nil {
		logrus.Error("Repository: Find menu error,", result.Error)
		return false, nil
	}

	if len(price) == len(menuNames) {
		return true, price
	}

	return false, nil
}
