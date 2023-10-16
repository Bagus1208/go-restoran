package handler

import (
	"fmt"
	"net/http"
	"restoran/features/order/model"
	"restoran/features/order/service"
	"restoran/helper"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderHandlerInterface interface {
	Insert() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Update() echo.HandlerFunc
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
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(fmt.Sprint("error when parshing data -", err.Error()), nil))
		}

		result, err := handler.service.Insert(orderInsert)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when inserting data -", err.Error()), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("successfully make an order", result))
	}
}

func (handler *orderHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := handler.service.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when getting all data -", err.Error()), nil))
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
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when getting data -", err.Error()), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get order by id", result))
	}
}

func (handler *orderHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		var orderUpdate model.OrderInput
		if err := c.Bind(&orderUpdate); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(fmt.Sprint("error when parshing data -", err.Error()), nil))
		}

		result, err := handler.service.Update(id, orderUpdate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when updating data -", err.Error()), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully updated data", result))
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
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when deleting data -", err.Error()), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully deleted data", nil))
	}
}
