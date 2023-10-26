package repository

import (
	"context"
	"errors"
	"mime/multipart"
	"restoran/config"
	"restoran/features/menu/model"
	"restoran/helper"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MenuRepositoryInterface interface {
	Insert(newData *model.Menu) (*model.Menu, error)
	GetAll(pagination model.QueryParam) ([]model.Menu, error)
	GetCategory(queryParam model.QueryParam) ([]model.Menu, error)
	GetFavorite() ([]model.Favorite, error)
	GetByName(name string) *model.Menu
	Update(id int, updateData *model.Menu) (*model.Menu, error)
	Delete(id int) error
	UploadImage(fileHeader *multipart.FileHeader, name string) (string, error)
	TotalData() (int64, error)
	TotalDataByCategory(category string) (int64, error)
}

type menuRepo struct {
	db     *gorm.DB
	cdn    *cloudinary.Cloudinary
	config config.Config
}

func NewMenuRepo(DB *gorm.DB, CDN *cloudinary.Cloudinary, config config.Config) MenuRepositoryInterface {
	return &menuRepo{
		db:     DB,
		cdn:    CDN,
		config: config,
	}
}

func (repository *menuRepo) Insert(newData *model.Menu) (*model.Menu, error) {
	result := repository.db.Create(newData)
	if result.Error != nil {
		logrus.Error("Repository: Inserting data error,", result.Error)
		return nil, result.Error
	}

	return newData, nil
}

func (repository *menuRepo) GetAll(pagination model.QueryParam) ([]model.Menu, error) {
	var menus []model.Menu
	var offset = (pagination.Page - 1) * pagination.PageSize

	result := repository.db.Offset(offset).Limit(pagination.PageSize).Find(&menus)
	if result.Error != nil {
		logrus.Error("Repository: Get all data error,", result.Error)
		return nil, result.Error
	}

	return menus, nil
}

func (repository *menuRepo) GetCategory(queryParam model.QueryParam) ([]model.Menu, error) {
	var menus []model.Menu
	var offset = (queryParam.Page - 1) * queryParam.PageSize

	result := repository.db.Where("category = ?", queryParam.Category).Offset(offset).Limit(queryParam.PageSize).Find(&menus)
	if result.Error != nil {
		logrus.Error("Repository: Get data by category error", result.Error)
		return nil, result.Error
	}

	if result == nil {
		return nil, errors.New("category not found")
	}

	return menus, nil
}

func (repository *menuRepo) GetByName(name string) *model.Menu {
	var menu model.Menu
	result := repository.db.Where("name = ?", name).First(&menu)
	if result.Error != nil {
		return nil
	}

	return &menu
}

func (repository *menuRepo) GetFavorite() ([]model.Favorite, error) {
	var favorites []model.Favorite
	result := repository.db.Table("order_details").
		Select("menu_name, SUM(quantity) AS total_order").
		Group("menu_name").
		Having("SUM(quantity) > ?", 20).
		Order("total_order DESC").
		Scan(&favorites)

	if result.Error != nil {
		logrus.Error("Repository: Get favorite data error", result.Error)
		return nil, result.Error
	}

	return favorites, nil
}

func (repository *menuRepo) Update(id int, updateData *model.Menu) (*model.Menu, error) {
	result := repository.db.Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		logrus.Error("Repository: Update data error,", result.Error)
		return nil, result.Error
	}

	var updatedUser = new(model.Menu)
	if err := repository.db.Where("id = ?", id).First(updatedUser).Error; err != nil {
		logrus.Error("Repository: Get data update error,", result.Error)
		return nil, err
	}

	return updatedUser, nil
}

func (repository *menuRepo) Delete(id int) error {
	var deleteMenu model.Menu
	deleteMenu.ID = uint(id)
	result := repository.db.Delete(&deleteMenu)
	if result.Error != nil {
		logrus.Error("Repository: Delete error,", result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		logrus.Error("Repository: Delete error,", result.Error)
		return errors.New("data not found")
	}

	return nil
}

func (repository *menuRepo) UploadImage(fileHeader *multipart.FileHeader, name string) (string, error) {
	var file = helper.OpenFileHeader(fileHeader)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := repository.cdn.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   repository.config.CDN_Folder_Name,
		PublicID: name,
	})
	if err != nil {
		logrus.Error("Repository: Upload image error,", err)
		return "", err
	}

	return response.SecureURL, nil
}

func (repository *menuRepo) TotalData() (int64, error) {
	var total int64

	result := repository.db.Table("menus").Where("deleted_at IS NULL").Count(&total)
	if result.Error != nil {
		return -1, result.Error
	}

	return total, nil
}

func (repository *menuRepo) TotalDataByCategory(category string) (int64, error) {
	var total int64

	result := repository.db.Table("menus").Where("category = ? AND deleted_at IS NULL", category).Count(&total)
	if result.Error != nil {
		return -1, result.Error
	}

	return total, nil
}
