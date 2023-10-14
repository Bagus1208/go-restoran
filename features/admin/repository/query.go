package repository

import (
	"restoran/features/admin"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminRepositoryInterface interface {
	Insert(newData *admin.Admin) (*admin.Admin, error)
	Login(email string, password string) (*admin.Admin, error)
}

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(DB *gorm.DB) AdminRepositoryInterface {
	return &adminRepo{
		db: DB,
	}
}

func (repository *adminRepo) Insert(newData *admin.Admin) (*admin.Admin, error) {
	result := repository.db.Create(newData)
	if result.Error != nil {
		logrus.Error("Repository: Inserting data error,", result.Error)
		return nil, result.Error
	}

	return newData, nil
}

func (repository *adminRepo) Login(email string, password string) (*admin.Admin, error) {
	var admin = new(admin.Admin)
	result := repository.db.Where("email =? and password =?", email, password).First(admin)
	if result.Error != nil {
		logrus.Error("Repository: Login error,", result.Error)
		return nil, result.Error
	}

	return admin, nil
}