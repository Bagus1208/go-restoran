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
	GetData() echo.HandlerFunc
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

func (handler *menuHandler) GetData() echo.HandlerFunc {
	return func(c echo.Context) error {
		var queryParam model.QueryParam

		queryParam.Page, _ = strconv.Atoi(c.QueryParam("page"))
		queryParam.PageSize, _ = strconv.Atoi(c.QueryParam("page_size"))
		queryParam.Name = c.QueryParam("name")
		queryParam.Category = c.QueryParam("category")

		if queryParam.Page < 1 || queryParam.PageSize < 1 {
			queryParam.Page = 1
			queryParam.PageSize = 10
		}

		var result []model.MenuResponse
		var paginationResponse *model.Pagination
		var err error
		var message string

		if queryParam.Name != "" {
			data, err := handler.service.GetByName(queryParam.Name)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
			}

			message := "successfully get menu: " + queryParam.Name

			return c.JSON(http.StatusOK, helper.FormatResponse(message, data))
		} else if queryParam.Category != "" {
			result, paginationResponse, err = handler.service.GetCategory(queryParam)
			message = "successfully get menu by category: " + queryParam.Category
		} else {
			result, paginationResponse, err = handler.service.GetAll(queryParam)
			message = "successfully get all menu"
		}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatPaginationResponse(message, result, paginationResponse))
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
