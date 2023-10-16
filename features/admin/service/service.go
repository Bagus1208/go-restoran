package service

import (
	"errors"
	"restoran/features/admin/model"
	"restoran/features/admin/repository"
	"restoran/helper"
	"strings"

	"github.com/go-playground/validator/v10"
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
	validator  *validator.Validate
}

func NewAdminService(repo repository.AdminRepositoryInterface, generate helper.GeneratorInterface, jwt helper.JWTInterface, validate *validator.Validate) AdminServiceInterface {
	return &adminService{
		repository: repo,
		generator:  generate,
		jwt:        jwt,
		validator:  validate,
	}
}

func (service *adminService) Insert(newData model.AdminInput) (*model.Admin, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var newUser = helper.RequestToAdmin(newData)

	newID, err := service.generator.GenerateUUID()
	if err != nil {
		logrus.Error("Service: Generating ID failed")
		return nil, errors.New("cannot generating id " + err.Error())
	}

	newUser.ID = newID
	result, err := service.repository.Insert(newUser)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("cannot insert data " + err.Error())
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
