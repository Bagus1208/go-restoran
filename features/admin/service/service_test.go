package service

import (
	"errors"
	"restoran/features/admin/mocks"
	"restoran/features/admin/model"
	"restoran/helper"
	mockHelper "restoran/helper/mocks"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	var jwt = mockHelper.NewJWTInterface(t)
	var generator = mockHelper.NewGeneratorInterface(t)
	var hash = mockHelper.NewHashInterface(t)
	var repository = mocks.NewAdminRepositoryInterface(t)
	var validate = validator.New()
	var service = NewAdminService(repository, jwt, generator, hash, validate)

	var newData = model.AdminInput{
		Name:     "Bagus",
		Email:    "bagus@gmail.com",
		Password: "bagus123",
	}

	var invalidData = model.AdminInput{
		Email:    "bagus@gmail.com",
		Password: "bagus123",
	}

	t.Run("success insert", func(t *testing.T) {
		newUser := helper.RequestToAdmin(newData)
		repository.On("IsDuplicateEmail", newData.Email).Return(false).Once()
		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		hash.On("HashPassword", newUser.Password).Return("hashPassword", nil).Once()

		newUser.ID = "randomUUID"
		newUser.Password = "hashPassword"
		repository.On("Insert", newUser).Return(newUser, nil).Once()

		result, err := service.Insert(newData)
		assert.Nil(t, err)
		assert.Equal(t, newUser.Name, result.Name)
		generator.AssertExpectations(t)
		hash.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		result, err := service.Insert(invalidData)

		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Email already used", func(t *testing.T) {
		repository.On("IsDuplicateEmail", newData.Email).Return(true).Once()

		result, err := service.Insert(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "email already used")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Generate id failed", func(t *testing.T) {
		repository.On("IsDuplicateEmail", newData.Email).Return(false).Once()
		generator.On("GenerateUUID").Return("", errors.New("id generator failed")).Once()

		result, err := service.Insert(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "id generator failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Hash password failed", func(t *testing.T) {
		repository.On("IsDuplicateEmail", newData.Email).Return(false).Once()
		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		hash.On("HashPassword", newData.Password).Return("", errors.New("hash password failed")).Once()

		result, err := service.Insert(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "hash password failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
		generator.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Insert data failed", func(t *testing.T) {
		newUser := helper.RequestToAdmin(newData)
		repository.On("IsDuplicateEmail", newData.Email).Return(false).Once()
		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		hash.On("HashPassword", newUser.Password).Return("hashPassword", nil).Once()

		newUser.ID = "randomUUID"
		newUser.Password = "hashPassword"
		repository.On("Insert", newUser).Return(nil, errors.New("cannot insert data")).Once()

		result, err := service.Insert(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "cannot insert data")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
		generator.AssertExpectations(t)
		hash.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	var jwt = mockHelper.NewJWTInterface(t)
	var generator = mockHelper.NewGeneratorInterface(t)
	var hash = mockHelper.NewHashInterface(t)
	var repository = mocks.NewAdminRepositoryInterface(t)
	var validate = validator.New()
	var service = NewAdminService(repository, jwt, generator, hash, validate)

	var adminData = model.Admin{
		ID:       "randomID",
		Name:     "Bagus",
		Email:    "bagus@gmail.com",
		Password: "bagus123",
	}

	var loginData = model.AdminInputLogin{
		Email:    "bagus@gmail.com",
		Password: "bagus123",
	}

	t.Run("success login", func(t *testing.T) {
		var jwtResult = map[string]any{"access_token": "randomAccessToken"}
		repository.On("Login", loginData.Email).Return(&adminData, nil).Once()
		hash.On("CompareHash", loginData.Password, adminData.Password).Return(true).Once()
		jwt.On("GenerateJWT", adminData.ID, adminData.Name, adminData.Email).Return(jwtResult).Once()

		result, err := service.Login(loginData.Email, loginData.Password)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Bagus", result.Name)
		assert.Equal(t, jwtResult, result.Access)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

	t.Run("login repository failed", func(t *testing.T) {
		repository.On("Login", loginData.Email).Return(nil, errors.New("data admin not found")).Once()

		result, err := service.Login(loginData.Email, loginData.Password)
		assert.Error(t, err)
		assert.EqualError(t, err, "email or passowrd is wrong")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("compare hash failed", func(t *testing.T) {
		repository.On("Login", loginData.Email).Return(&adminData, nil).Once()
		hash.On("CompareHash", loginData.Password, adminData.Password).Return(false).Once()

		result, err := service.Login(loginData.Email, loginData.Password)
		assert.Error(t, err)
		assert.EqualError(t, err, "email or passowrd is wrong")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("get access token failed", func(t *testing.T) {
		repository.On("Login", loginData.Email).Return(&adminData, nil).Once()
		hash.On("CompareHash", loginData.Password, adminData.Password).Return(true).Once()
		jwt.On("GenerateJWT", adminData.ID, adminData.Name, adminData.Email).Return(nil).Once()

		result, err := service.Login(loginData.Email, loginData.Password)
		assert.Error(t, err)
		assert.EqualError(t, err, "get token process failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})
}

func TestSetNoTable(t *testing.T) {
	var jwt = mockHelper.NewJWTInterface(t)
	var generator = mockHelper.NewGeneratorInterface(t)
	var hash = mockHelper.NewHashInterface(t)
	var repository = mocks.NewAdminRepositoryInterface(t)
	var validate = validator.New()
	var service = NewAdminService(repository, jwt, generator, hash, validate)

	var setTable = model.InputTable{
		TableNumber: 1,
	}

	var invalidTable = model.InputTable{
		TableNumber: 0,
	}

	adminName := "Bagus"

	t.Run("success set no table", func(t *testing.T) {
		var jwtResult = "randomAccessToken"
		jwt.On("GenerateTableToken", setTable.TableNumber, adminName).Return(jwtResult).Once()

		result, err := service.SetNoTable(adminName, setTable)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, jwtResult, result)
		jwt.AssertExpectations(t)
	})

	t.Run("validation failed", func(t *testing.T) {
		result, err := service.SetNoTable(adminName, invalidTable)
		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Empty(t, result)
	})

	t.Run("get token failed", func(t *testing.T) {
		jwt.On("GenerateTableToken", setTable.TableNumber, adminName).Return("").Once()

		result, err := service.SetNoTable(adminName, setTable)
		assert.Error(t, err)
		assert.EqualError(t, err, "get token process failed")
		assert.Empty(t, result)
		jwt.AssertExpectations(t)
	})
}
