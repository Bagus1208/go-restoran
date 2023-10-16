package handler

import (
	"net/http"
	"restoran/features/admin/model"
	"restoran/features/admin/service"
	"restoran/helper"
	"strings"

	"github.com/labstack/echo/v4"
)

type AdminHandlerInterface interface {
	Insert() echo.HandlerFunc
	Login() echo.HandlerFunc
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
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully insert data", result))
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
				return c.JSON(http.StatusNotFound, helper.FormatResponse("user not found", nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully login", result))
	}
}
