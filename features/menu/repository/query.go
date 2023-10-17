package repository

import (
	"context"
	"mime/multipart"
	"os"
	"restoran/features/menu/model"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MenuRepositoryInterface interface {
	Insert(newData *model.Menu) (*model.Menu, error)
	GetAll(pagination model.Pagination) ([]model.Menu, error)
	GetCategory(category string, pagination model.Pagination) ([]model.Menu, error)
	GetByName(name string) *model.Menu
	Update(id int, updateData *model.Menu) (*model.Menu, error)
	Delete(id int) error
	UploadImage(ctx context.Context, file multipart.File, name string) (string, error)
}

type menuRepo struct {
	db  *gorm.DB
	cdn *cloudinary.Cloudinary
}

func NewMenuRepo(DB *gorm.DB, CDN *cloudinary.Cloudinary) MenuRepositoryInterface {
	return &menuRepo{
		db:  DB,
		cdn: CDN,
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

func (repository *menuRepo) GetAll(pagination model.Pagination) ([]model.Menu, error) {
	var menus []model.Menu
	var offset = (pagination.Page - 1) * pagination.PageSize

	result := repository.db.Offset(offset).Limit(pagination.PageSize).Find(&menus)
	if result.Error != nil {
		logrus.Error("Repository: Get all data error,", result.Error)
		return nil, result.Error
	}

	return menus, nil
}

func (repository *menuRepo) GetCategory(category string, pagination model.Pagination) ([]model.Menu, error) {
	var menus []model.Menu
	var offset = (pagination.Page - 1) * pagination.PageSize

	result := repository.db.Where("category = ?", category).Offset(offset).Limit(pagination.PageSize).Find(&menus)
	if result.Error != nil {
		logrus.Error("Repository: Get data by category error", result.Error)
		return nil, result.Error
	}

	return menus, nil
}

func (repository *menuRepo) GetByName(name string) *model.Menu {
	var menu model.Menu
	result := repository.db.Where("name =?", name).First(&menu)
	if result.Error != nil {
		return nil
	}

	return &menu
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

	return nil
}

func (repository *menuRepo) UploadImage(ctx context.Context, file multipart.File, name string) (string, error) {
	response, err := repository.cdn.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   os.Getenv("CLOUDINARY_UPLOAD_FOLDER_NAME"),
		PublicID: name,
	})
	if err != nil {
		logrus.Error("Repository: Upload image error,", err)
		return "", err
	}

	return response.SecureURL, nil
}
