package service

import (
	"errors"
	"restoran/features/admin"
	"restoran/features/admin/model"
	"restoran/features/admin/repository"
	"restoran/helper"
	"strings"

	"github.com/sirupsen/logrus"
)

type AdminServiceInterface interface {
	Insert(newData model.AdminInput) (*admin.Admin, error)
	Login(email string, password string) (*admin.UserCredential, error)
}

type adminService struct {
	repository repository.AdminRepositoryInterface
	generator  helper.GeneratorInterface
	jwt        helper.JWTInterface
}

func NewAdminService(repo repository.AdminRepositoryInterface, generate helper.GeneratorInterface, jwt helper.JWTInterface) AdminServiceInterface {
	return &adminService{
		repository: repo,
		generator:  generate,
		jwt:        jwt,
	}
}

func (service *adminService) Insert(newData model.AdminInput) (*admin.Admin, error) {
	var newUser = helper.RequestToAdmin(newData)
	result, err := service.repository.Insert(newUser)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("Cannot insert data " + err.Error())
	}

	return result, nil
}

func (service *adminService) Login(email string, password string) (*admin.UserCredential, error) {
	result, err := service.repository.Login(email, password)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.New("data not found")
		}
		return nil, errors.New("process failed")
	}

	token := service.jwt.GenerateJWT(result.ID)
	if token == nil {
		return nil, errors.New("get token process failed")
	}

	var response = new(admin.UserCredential)
	response.Name = result.Name
	response.Access = token

	return response, nil
}
