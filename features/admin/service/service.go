package service

import (
	"errors"
	"restoran/features/admin/model"
	"restoran/features/admin/repository"
	"restoran/helper"
	"strings"

	"github.com/sirupsen/logrus"
)

type AdminServiceInterface interface {
	Insert(newData model.AdminInput) (*model.Admin, error)
	Login(email string, password string) (*model.UserCredential, error)
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

func (service *adminService) Insert(newData model.AdminInput) (*model.Admin, error) {
	var newUser = helper.RequestToAdmin(newData)

	newID, err := service.generator.GenerateUUID()
	if err != nil {
		logrus.Error("Service: Generating ID failed")
		return nil, errors.New("Cannot generating ID " + err.Error())
	}

	newUser.ID = newID
	result, err := service.repository.Insert(newUser)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("Cannot insert data " + err.Error())
	}

	return result, nil
}

func (service *adminService) Login(email string, password string) (*model.UserCredential, error) {
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

	var response = new(model.UserCredential)
	response.Name = result.Name
	response.Access = token

	return response, nil
}
