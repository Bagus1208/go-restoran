package handler

import (
	"net/http"
	"restoran/features/admin/model"
	"restoran/features/admin/service"
	"restoran/helper"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AdminHandlerInterface interface {
	Insert() echo.HandlerFunc
	Login() echo.HandlerFunc
	SetNoTable() echo.HandlerFunc
}

type adminHandler struct {
	service service.AdminServiceInterface
}

func NewAdminHandler(service service.AdminServiceInterface) AdminHandlerInterface {
	return &adminHandler{
		service: service,
	}
}

func (handler *adminHandler) Insert() echo.HandlerFunc {
	return func(c echo.Context) error {
		var adminInsert model.AdminInput
		if err := c.Bind(&adminInsert); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		result, err := handler.service.Insert(adminInsert)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(err.Error(), nil))
			} else if strings.Contains(err.Error(), "email already used") {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(err.Error(), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("successfully insert data", result))
	}
}

func (handler *adminHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var adminLogin model.AdminInputLogin
		if err := c.Bind(&adminLogin); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		result, err := handler.service.Login(adminLogin.Email, adminLogin.Password)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, helper.FormatResponse(err.Error(), nil))
			} else if strings.Contains(err.Error(), "email or passowrd is wrong") {
				return c.JSON((http.StatusBadRequest), helper.FormatResponse(err.Error(), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully login", result))
	}
}

func (handler *adminHandler) SetNoTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		var setTable model.InputTable
		if err := c.Bind(&setTable); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		var adminName string
		if user, ok := c.Get("user").(*jwt.Token); ok {
			if claims, ok := user.Claims.(jwt.MapClaims); ok {
				if name, ok := claims["name"].(string); ok {
					adminName = name
				}
			}
		}

		result, err := handler.service.SetNoTable(adminName, setTable)
		if err != nil {
			if strings.Contains(err.Error(), "validation") {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(err.Error(), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get token table", result))
	}
}
