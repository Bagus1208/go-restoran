package handler

import (
	"net/http"
	"restoran/features/order/model"
	"restoran/features/order/service"
	"restoran/helper"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type OrderHandlerInterface interface {
	Insert() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type orderHandler struct {
	service service.OrderServiceInterface
}

func NewOrderHandler(service service.OrderServiceInterface) OrderHandlerInterface {
	return &orderHandler{
		service: service,
	}
}

func (handler *orderHandler) Insert() echo.HandlerFunc {
	return func(c echo.Context) error {
		var orderInsert model.OrderInput
		if err := c.Bind(&orderInsert); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		var stringToken = c.Request().Header.Get("Authorization")
		if stringToken == "" {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("token not found", nil))
		}

		noTable, err := helper.ExtractToken(stringToken)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(err.Error(), nil))
		}

		orderInsert.NoTable = noTable
		result, err := handler.service.Insert(orderInsert)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("successfully make an order", result))
	}
}

func (handler *orderHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		var pagination model.Pagination

		pagination.Page, _ = strconv.Atoi(c.QueryParam("page"))
		pagination.PageSize, _ = strconv.Atoi(c.QueryParam("page_size"))

		if pagination.Page < 1 || pagination.PageSize < 1 {
			pagination.Page = 1
			pagination.PageSize = 10
		}

		result, err := handler.service.GetAll(pagination)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get all orders", result))
	}
}

func (handler *orderHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		result, err := handler.service.GetByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get order by id", result))
	}
}

func (handler *orderHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		err = handler.service.Delete(id)
		if err != nil {
			if strings.Contains(err.Error(), "no rows affected") {
				return c.JSON(http.StatusNotFound, helper.FormatResponse(err.Error(), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully deleted data", nil))
	}
}
