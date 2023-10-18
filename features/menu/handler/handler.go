package handler

import (
	"net/http"
	"restoran/features/menu/model"
	"restoran/features/menu/service"
	"restoran/helper"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type MenuHandlerInterface interface {
	Insert() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetCategory() echo.HandlerFunc
	GetFavorite() echo.HandlerFunc
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

func (handler *menuHandler) Insert() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, helper.FormatResponse("image not found", nil))
		}

		var menuInsert model.MenuInput
		if err := c.Bind(&menuInsert); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		result, err := handler.service.Insert(fileHeader, menuInsert)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(err.Error(), nil))
			}

			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("successfully insert data", result))
	}
}

func (handler *menuHandler) GetAll() echo.HandlerFunc {
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

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get all menu", result))
	}
}

func (handler *menuHandler) GetCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		category := c.Param("category")
		if category == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("category is required", nil))
		}

		var pagination model.Pagination

		pagination.Page, _ = strconv.Atoi(c.QueryParam("page"))
		pagination.PageSize, _ = strconv.Atoi(c.QueryParam("page_size"))

		if pagination.Page < 1 || pagination.PageSize < 1 {
			pagination.Page = 1
			pagination.PageSize = 10
		}

		result, err := handler.service.GetCategory(category, pagination)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get menu by category", result))
	}
}

func (handler *menuHandler) GetFavorite() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := handler.service.GetFavorite()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully get favorite menu", result))
	}
}

func (handler *menuHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("id is required", nil))
		}

		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, helper.FormatResponse("image not found", nil))
		}

		var menuUpdate model.MenuInput
		if err := c.Bind(&menuUpdate); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("error when parshing data", nil))
		}

		result, err := handler.service.Update(id, fileHeader, menuUpdate)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(err.Error(), nil))
			}

			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully update data", result))
	}
}

func (handler *menuHandler) Delete() echo.HandlerFunc {
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

		return c.JSON(http.StatusOK, helper.FormatResponse("successfully delete data", nil))
	}
}
