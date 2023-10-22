package handler

import (
	"encoding/json"
	"net/http"
	"restoran/features/transaction/model"
	"restoran/features/transaction/service"
	"restoran/helper"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type TransactionHandlerInterface interface {
	Insert() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Notifications() echo.HandlerFunc
}

type transactionHandler struct {
	service service.TransactionServiceInterface
}

func NewTransactionHandler(service service.TransactionServiceInterface) TransactionHandlerInterface {
	return &transactionHandler{
		service: service,
	}
}

func (handler *transactionHandler) Insert() echo.HandlerFunc {
	return func(c echo.Context) error {
		var transaction model.TransactionInput
		if err := c.Bind(&transaction); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		result, err := handler.service.Insert(transaction)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(err.Error(), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("successfully create transaction", result))
	}
}

func (handler *transactionHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		var pagination model.QueryParam

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

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get all transaction", result))
	}
}

func (handler *transactionHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		result, err := handler.service.GetByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get transaction by id", result))
	}
}

func (handler *transactionHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		err = handler.service.Delete(id)
		if err != nil {
			if strings.Contains(err.Error(), "data not found") {
				return c.JSON(http.StatusNotFound, helper.FormatResponse(err.Error(), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully delete data", nil))
	}
}

func (handler *transactionHandler) Notifications() echo.HandlerFunc {
	return func(c echo.Context) error {
		var notificationPayload map[string]any

		if err := json.NewDecoder(c.Request().Body).Decode(&notificationPayload); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		err := handler.service.Notifications(notificationPayload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
	}
}
