package repository

import (
	"restoran/features/admin/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminRepositoryInterface interface {
	Insert(newData *model.Admin) (*model.Admin, error)
	Login(email string) (*model.Admin, error)
	IsDuplicateEmail(email string) bool
}

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(DB *gorm.DB) AdminRepositoryInterface {
	return &adminRepo{
		db: DB,
	}
}

func (repository *adminRepo) Insert(newData *model.Admin) (*model.Admin, error) {
	result := repository.db.Create(newData)
	if result.Error != nil {
		logrus.Error("Repository: Inserting data error,", result.Error)
		return nil, result.Error
	}

	return newData, nil
}

func (repository *adminRepo) Login(email string) (*model.Admin, error) {
	var admin = new(model.Admin)
	result := repository.db.Where("email = ?", email).First(admin)
	if result.Error != nil {
		logrus.Error("Repository: Login error,", result.Error)
		return nil, result.Error
	}

	return admin, nil
}

func (repository *adminRepo) IsDuplicateEmail(email string) bool {
	var count int64
	result := repository.db.Model(&model.Admin{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		logrus.Error("Repository: IsDuplicateEmail error,", result.Error)
		return true 
	}

	return count > 0
}


