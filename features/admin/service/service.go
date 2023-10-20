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
	Insert(newData model.AdminInput) (*model.AdminResponse, error)
	Login(email string, password string) (*model.UserCredential, error)
	SetNoTable(noTable int, email string, password string) (string, error)
}

type adminService struct {
	repository repository.AdminRepositoryInterface
	jwt        helper.JWTInterface
	validator  *validator.Validate
}

func NewAdminService(repo repository.AdminRepositoryInterface, jwt helper.JWTInterface, validate *validator.Validate) AdminServiceInterface {
	return &adminService{
		repository: repo,
		jwt:        jwt,
		validator:  validate,
	}
}

func (service *adminService) Insert(newData model.AdminInput) (*model.AdminResponse, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	var newUser = helper.RequestToAdmin(newData)

	newUser.ID = helper.GenerateUUID()
	newUser.Password = helper.HashPassword(newUser.Password)

	result, err := service.repository.Insert(newUser)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("cannot insert data: " + err.Error())
	}

	var adminResponse = helper.AdminToResponse(result)

	return adminResponse, nil
}

func (service *adminService) Login(email string, password string) (*model.UserCredential, error) {
	result, err := service.repository.Login(email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.New("user admin not found")
		}
		return nil, errors.New("process failed")
	}

	if !helper.CompareHash(password, result.Password) {
		return nil, errors.New("wrong password")
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

func (service *adminService) SetNoTable(noTable int, email string, password string) (string, error) {
	result, err := service.repository.Login(email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", errors.New("user admin not found")
		}
		return "", errors.New("process failed")
	}

	if !helper.CompareHash(password, result.Password) {
		return "", errors.New("wrong password")
	}

	token := service.jwt.GenerateTableToken(noTable, result.Name)
	if token == "" {
		return "", errors.New("get token process failed")
	}

	return token, nil
}
