package handler

import (
	"fmt"
	"net/http"
	"restoran/features/menu/model"
	"restoran/features/menu/service"
	"restoran/helper"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MenuHandlerInterface interface {
	Insert() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetCategory() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type menuHandler struct {
	service service.MenuServiceInterface
}

func NewMenuHandler(service service.MenuServiceInterface) MenuHandlerInterface {
	return &menuHandler{
		service: service,
	}
}

func (mh *menuHandler) Insert() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, helper.FormatResponse("unprocessable content -", err.Error()))
		}

		var menuInsert model.MenuInput
		if err := c.Bind(&menuInsert); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(fmt.Sprint("error when parshing data -", err.Error()), nil))
		}

		result, err := mh.service.Insert(fileHeader, menuInsert)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when inserting data -", err.Error()), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("successfully inserted data", result))
	}
}

func (mh *menuHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := mh.service.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when getting data -", err.Error()), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get all menu", result))
	}
}

func (mh *menuHandler) GetCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		category := c.Param("category")
		if category == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("category is required", nil))
		}

		result, err := mh.service.GetCategory(category)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when getting data -", err.Error()), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get menu by category", result))
	}
}

func (mh *menuHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, helper.FormatResponse("unprocessable content -", err.Error()))
		}

		var menuUpdate model.MenuInput
		if err := c.Bind(&menuUpdate); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(fmt.Sprint("error when parshing data -", err.Error()), nil))
		}

		result, err := mh.service.Update(id, fileHeader, menuUpdate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when updating data -", err.Error()), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully updated data", result))
	}
}

func (mh *menuHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		err = mh.service.Delete(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(fmt.Sprint("error when deleting data -", err.Error()), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully deleted data", nil))
	}
}
