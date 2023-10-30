package service

import (
	"errors"
	"mime/multipart"
	"restoran/features/menu/mocks"
	"restoran/features/menu/model"
	"restoran/helper"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	var newData = model.MenuInput{
		Name:        "Martabak",
		Category:    "Makanan",
		Price:       25000,
		Description: "Makanan manis yang nikmat",
		Image:       "www.cloudinary.com/images/martabak",
	}

	var invalidData = model.MenuInput{
		Name: "Martabak",
	}

	var menu = model.Menu{
		ID:          1,
		Name:        "Martabak",
		Category:    "Makanan",
		Price:       25000,
		Description: "Makanan manis yang nikmat",
		Image:       "www.cloudinary.com/images/martabak",
	}

	var fileHeader *multipart.FileHeader

	t.Run("success insert menu", func(t *testing.T) {
		var newMenu = helper.RequestToMenu(newData)

		repository.On("GetByName", newData.Name).Return(nil).Once()
		repository.On("UploadImage", fileHeader, newData.Name).Return(menu.Image, nil).Once()
		repository.On("Insert", newMenu).Return(&menu, nil).Once()

		result, err := service.Insert(fileHeader, newData)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newData.Name, result.Name)
		repository.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		result, err := service.Insert(fileHeader, invalidData)
		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Menu already exists", func(t *testing.T) {
		repository.On("GetByName", newData.Name).Return(&menu).Once()

		result, err := service.Insert(fileHeader, newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "menu already exists")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Upload image failed", func(t *testing.T) {
		repository.On("GetByName", newData.Name).Return(nil).Once()
		repository.On("UploadImage", fileHeader, newData.Name).Return("", errors.New("upload image error")).Once()

		result, err := service.Insert(fileHeader, newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "upload image failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Insert data failed", func(t *testing.T) {
		var newMenu = helper.RequestToMenu(newData)

		repository.On("GetByName", newData.Name).Return(nil).Once()
		repository.On("UploadImage", fileHeader, newData.Name).Return(menu.Image, nil).Once()
		repository.On("Insert", newMenu).Return(nil, errors.New("insert data error")).Once()

		result, err := service.Insert(fileHeader, newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "insert data menu failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	var pagination = model.QueryParam{
		Page:     1,
		PageSize: 2,
	}

	var paginationDetail = model.Pagination{
		Page:       1,
		PageSize:   2,
		TotalItems: 10,
	}

	var listMenu = []model.Menu{
		{
			ID:          1,
			Name:        "Martabak",
			Category:    "Makanan",
			Description: "Makanan manis yang nikmat",
			Price:       20000,
			Image:       "www.cloudinary.com/images/martabak",
		},
		{
			ID:          2,
			Name:        "Jus alpukat",
			Category:    "Minuman",
			Description: "Alpukat diblender dan diberi susu kental manis",
			Price:       10000,
			Image:       "www.cloudinary.com/images/jus_alpukat",
		},
	}

	t.Run("Success get data", func(t *testing.T) {
		repository.On("GetAll", pagination).Return(listMenu, nil).Once()
		repository.On("TotalData").Return(paginationDetail.TotalItems, nil).Once()

		result, paginationResponse, err := service.GetAll(pagination)
		assert.Nil(t, err)
		assert.Equal(t, len(listMenu), len(result))
		assert.Equal(t, listMenu[0].Name, result[0].Name)
		assert.Equal(t, paginationDetail.TotalItems, paginationResponse.TotalItems)
		repository.AssertExpectations(t)
	})

	t.Run("Get data failed", func(t *testing.T) {
		repository.On("GetAll", pagination).Return(nil, errors.New("get data error")).Once()

		result, paginationResponse, err := service.GetAll(pagination)
		assert.Error(t, err)
		assert.EqualError(t, err, "get data menu failed")
		assert.Nil(t, result)
		assert.Nil(t, paginationResponse)
		repository.AssertExpectations(t)
	})

	t.Run("Get total data failed", func(t *testing.T) {
		repository.On("GetAll", pagination).Return(listMenu, nil).Once()
		repository.On("TotalData").Return(int64(-1), errors.New("get total data errors")).Once()

		result, paginationResponse, err := service.GetAll(pagination)
		assert.Error(t, err)
		assert.EqualError(t, err, "get total menu failed")
		assert.Nil(t, result)
		assert.Nil(t, paginationResponse)
		repository.AssertExpectations(t)

	})

}

func TestGetCategory(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	var pagination = model.QueryParam{
		Page:     1,
		PageSize: 2,
		Category: "makanan",
	}

	var paginationDetail = model.Pagination{
		Page:       1,
		PageSize:   2,
		TotalItems: 10,
	}

	var listMenu = []model.Menu{
		{
			ID:          1,
			Name:        "Martabak",
			Category:    "Makanan",
			Description: "Makanan manis yang nikmat",
			Price:       20000,
			Image:       "www.cloudinary.com/images/martabak",
		},
		{
			ID:          2,
			Name:        "Mie goreng",
			Category:    "Makanan",
			Description: "Mie goreng dengan berbagai topping",
			Price:       10000,
			Image:       "www.cloudinary.com/images/jus_alpukat",
		},
	}

	t.Run("Succcess get data by category", func(t *testing.T) {
		repository.On("GetCategory", pagination).Return(listMenu, nil).Once()
		repository.On("TotalDataByCategory", pagination.Category).Return(paginationDetail.TotalItems, nil).Once()

		result, paginationResponse, err := service.GetCategory(pagination)
		assert.Nil(t, err)
		assert.Equal(t, len(listMenu), len(result))
		assert.Equal(t, listMenu[0].Category, result[0].Category)
		assert.Equal(t, paginationDetail.TotalItems, paginationResponse.TotalItems)
		repository.AssertExpectations(t)
	})

	t.Run("Get data by category failed", func(t *testing.T) {
		repository.On("GetCategory", pagination).Return(nil, errors.New("get data error")).Once()

		result, paginationResponse, err := service.GetCategory(pagination)
		assert.Error(t, err)
		assert.EqualError(t, err, "get data menu by category failed")
		assert.Nil(t, result)
		assert.Nil(t, paginationResponse)
		repository.AssertExpectations(t)
	})

	t.Run("get total data by category failed", func(t *testing.T) {
		repository.On("GetCategory", pagination).Return(listMenu, nil).Once()
		repository.On("TotalDataByCategory", pagination.Category).Return(int64(-1), errors.New("get total error")).Once()

		result, paginationResponse, err := service.GetCategory(pagination)
		assert.Error(t, err)
		assert.EqualError(t, err, "get total menu by category failed")
		assert.Nil(t, result)
		assert.Nil(t, paginationResponse)
		repository.AssertExpectations(t)
	})
}

func TestGetFavorite(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	var favoriteMenu = []model.Favorite{
		{
			MenuName:   "Martabak",
			TotalOrder: 100,
		},
	}

	t.Run("Success get data favorite", func(t *testing.T) {
		repository.On("GetFavorite").Return(favoriteMenu, nil).Once()

		result, err := service.GetFavorite()
		assert.Nil(t, err)
		assert.Equal(t, favoriteMenu[0].MenuName, result[0].MenuName)
		repository.AssertExpectations(t)
	})

	t.Run("Get data favorite failed", func(t *testing.T) {
		repository.On("GetFavorite").Return(nil, errors.New("get data error")).Once()

		result, err := service.GetFavorite()
		assert.Error(t, err)
		assert.EqualError(t, err, "get data favorite menu failed")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestGetByName(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	var menu = model.Menu{
		ID:          1,
		Name:        "Martabak",
		Category:    "Makanan",
		Price:       25000,
		Description: "Makanan manis yang nikmat",
		Image:       "www.cloudinary.com/images/martabak",
	}

	t.Run("Success get by name", func(t *testing.T) {
		var menuName = "Martabak"
		repository.On("GetByName", menuName).Return(&menu).Once()

		result, err := service.GetByName(menuName)
		assert.Nil(t, err)
		assert.Equal(t, menuName, result.Name)
		repository.AssertExpectations(t)
	})

	t.Run("Get by name failed", func(t *testing.T) {
		var menuName = "Martabak"
		repository.On("GetByName", menuName).Return(nil).Once()

		result, err := service.GetByName(menuName)
		assert.Error(t, err)
		assert.EqualError(t, err, "menu not found")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	var updateData = model.MenuInput{
		Name:        "Martabak",
		Category:    "Makanan",
		Price:       25000,
		Description: "Makanan manis yang nikmat",
		Image:       "www.cloudinary.com/images/martabak",
	}

	var invalidData = model.MenuInput{
		Name: "Martabak",
	}

	var menu = model.Menu{
		ID:          1,
		Name:        "Martabak",
		Category:    "Makanan",
		Price:       25000,
		Description: "Makanan manis yang nikmat",
		Image:       "www.cloudinary.com/images/martabak",
	}

	var fileHeader *multipart.FileHeader

	t.Run("Success update data", func(t *testing.T) {
		var updateMenu = helper.RequestToMenu(updateData)

		repository.On("UploadImage", fileHeader, updateData.Name).Return(menu.Image, nil).Once()
		repository.On("Update", 1, updateMenu).Return(&menu, nil).Once()

		result, err := service.Update(1, fileHeader, updateData)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updateData.Name, result.Name)
		repository.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		result, err := service.Update(1, fileHeader, invalidData)
		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Upload image failed", func(t *testing.T) {
		repository.On("UploadImage", fileHeader, updateData.Name).Return("", errors.New("upload image error")).Once()

		result, err := service.Update(1, fileHeader, updateData)
		assert.Error(t, err)
		assert.EqualError(t, err, "cannot upload image")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Update data failed", func(t *testing.T) {
		var updateMenu = helper.RequestToMenu(updateData)

		repository.On("UploadImage", fileHeader, updateData.Name).Return(menu.Image, nil).Once()
		repository.On("Update", 1, updateMenu).Return(nil, errors.New("update data error")).Once()

		result, err := service.Update(1, fileHeader, updateData)
		assert.Error(t, err)
		assert.EqualError(t, err, "cannot update data")
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	t.Run("success delete data", func(t *testing.T) {
		repository.On("Delete", 1).Return(nil).Once()

		err := service.Delete(1)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("delete data failed", func(t *testing.T) {
		repository.On("Delete", 1).Return(errors.New("delete data error")).Once()

		err := service.Delete(1)
		assert.Error(t, err)
		assert.EqualError(t, err, "delete data error")
		repository.AssertExpectations(t)
	})
}

func TestRecommendationMenu(t *testing.T) {
	var repository = mocks.NewMenuRepositoryInterface(t)
	var validate = validator.New()
	var service = NewMenuService(repository, validate)

	var request = model.RecommendationRequest{
		Message: "Saya ingin makanan yang pedas",
	}

	var invalidRequest = model.RecommendationRequest{}

	var menuName = []string{
		"Bakso",
		"Nasi goreng",
		"Kopi",
		"klepon",
	}

	t.Run("Success recommendation menu", func(t *testing.T) {
		request.MenuName = menuName

		repository.On("GetAllMenuName").Return(menuName, nil).Once()
		repository.On("RecommendationMenu", request).Return("Saya merekomendasikan bakso", nil).Once()

		result, err := service.RecommendationMenu(request)
		assert.Nil(t, err)
		assert.Equal(t, "Saya merekomendasikan bakso", result)
		repository.AssertExpectations(t)
	})

	t.Run("Validation failed", func(t *testing.T) {
		result, err := service.RecommendationMenu(invalidRequest)

		assert.Error(t, err)
		assert.EqualError(t, err, "validation failed please check your input and try again")
		assert.Empty(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Get menu failed", func(t *testing.T) {
		repository.On("GetAllMenuName").Return(nil, errors.New("get all menu name error")).Once()

		result, err := service.RecommendationMenu(request)
		assert.Error(t, err)
		assert.EqualError(t, err, "get menu name failed")
		assert.Empty(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Recommendation menu failed", func(t *testing.T) {
		request.MenuName = menuName

		repository.On("GetAllMenuName").Return(menuName, nil).Once()
		repository.On("RecommendationMenu", request).Return("", errors.New("recommendation error")).Once()

		result, err := service.RecommendationMenu(request)
		assert.Error(t, err)
		assert.EqualError(t, err, "recommendation error")
		assert.Empty(t, result)
		repository.AssertExpectations(t)
	})
}
