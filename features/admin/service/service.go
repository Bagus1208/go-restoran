package service

import (
	"errors"
	"restoran/features/admin/model"
	"restoran/features/admin/repository"
	"restoran/helper"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AdminServiceInterface interface {
	Insert(newData model.AdminInput) (*model.AdminResponse, error)
	Login(email string, password string) (*model.UserCredential, error)
	SetNoTable(adminName string, setTable model.InputTable) (string, error)
}

type adminService struct {
	repository repository.AdminRepositoryInterface
	jwt        helper.JWTInterface
	generate   helper.GeneratorInterface
	hash       helper.HashInterface
	validator  *validator.Validate
}

func NewAdminService(repo repository.AdminRepositoryInterface,
	jwt helper.JWTInterface,
	generator helper.GeneratorInterface,
	hashPassword helper.HashInterface,
	validate *validator.Validate) AdminServiceInterface {
	return &adminService{
		repository: repo,
		jwt:        jwt,
		generate:   generator,
		hash:       hashPassword,
		validator:  validate,
	}
}

func (service *adminService) Insert(newData model.AdminInput) (*model.AdminResponse, error) {
	err := service.validator.Struct(newData)
	if err != nil {
		return nil, errors.New("validation failed please check your input and try again")
	}

	if service.repository.IsDuplicateEmail(newData.Email) {
		return nil, errors.New("email already used")
	}

	var newUser = helper.RequestToAdmin(newData)

	newUser.ID, err = service.generate.GenerateUUID()
	if err != nil {
		return nil, errors.New("id generator failed")
	}

	newUser.Password, err = service.hash.HashPassword(newUser.Password)
	if err != nil {
		return nil, errors.New("hash password failed")
	}

	result, err := service.repository.Insert(newUser)
	if err != nil {
		logrus.Error("Service: Insert data failed,", err)
		return nil, errors.New("cannot insert data")
	}

	var adminResponse = helper.AdminToResponse(result)

	return adminResponse, nil
}

func (service *adminService) Login(email string, password string) (*model.UserCredential, error) {
	result, err := service.repository.Login(email)
	if err != nil {
		return nil, errors.New("email or passowrd is wrong")
	}

	if !service.hash.CompareHash(password, result.Password) {
		return nil, errors.New("email or passowrd is wrong")
	}

	token := service.jwt.GenerateJWT(result.ID, result.Name, result.Email)
	if token == nil {
		return nil, errors.New("get token process failed")
	}

	var response = helper.AdminToLoginResponse(result, token)

	return response, nil
}

func (service *adminService) SetNoTable(adminName string, setTable model.InputTable) (string, error) {
	if setTable.TableNumber == 0 && setTable.NoTable != 0 {
		setTable.TableNumber = setTable.NoTable
	} else if setTable.NoTable == 0 && setTable.TableNumber != 0 {
		setTable.NoTable = setTable.TableNumber
	}

	err := service.validator.Struct(setTable)
	if err != nil {
		return "", errors.New("validation failed please check your input and try again")
	}

	token := service.jwt.GenerateTableToken(setTable.TableNumber, adminName)
	if token == "" {
		return "", errors.New("get token process failed")
	}

	return token, nil
}
