package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"restoran/config"
	"restoran/features/menu/model"
	"restoran/helper"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MenuRepositoryInterface interface {
	Insert(newData *model.Menu) (*model.Menu, error)
	GetAll(pagination model.QueryParam) ([]model.Menu, error)
	GetCategory(queryParam model.QueryParam) ([]model.Menu, error)
	GetFavorite() ([]model.Favorite, error)
	GetByName(name string) *model.Menu
	GetAllMenuName() ([]string, error)
	Update(id int, updateData *model.Menu) (*model.Menu, error)
	Delete(id int) error
	UploadImage(fileHeader *multipart.FileHeader, name string) (string, error)
	TotalData() (int64, error)
	TotalDataByCategory(category string) (int64, error)
	RecommendationMenu(request model.RecommendationRequest) (string, error)
}

type menuRepo struct {
	db     *gorm.DB
	cdn    *cloudinary.Cloudinary
	ai     *openai.Client
	config config.Config
}

func NewMenuRepo(DB *gorm.DB, CDN *cloudinary.Cloudinary, AI *openai.Client, config config.Config) MenuRepositoryInterface {
	return &menuRepo{
		db:     DB,
		cdn:    CDN,
		ai:     AI,
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
	result := repository.db.Table("order_details AS od").
		Select("menus.name AS menu_name, SUM(od.quantity) AS total_order").
		Joins("JOIN menus ON menus.id = od.menu_id").
		Group("menus.name").
		Having("SUM(od.quantity) > 50").
		Order("total_order DESC").
		Scan(&favorites)

	if result.Error != nil {
		logrus.Error("Repository: Get favorite data error", result.Error)
		return nil, result.Error
	}

	return favorites, nil
}

func (repository *menuRepo) GetAllMenuName() ([]string, error) {
	var menuName []string

	result := repository.db.Table("menus").Select("name").Pluck("name", &menuName)
	if result.Error != nil {
		return nil, result.Error
	}

	return menuName, nil
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

func (repository *menuRepo) RecommendationMenu(request model.RecommendationRequest) (string, error) {
	ctx := context.TODO()

	chatMessage := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		Content: fmt.Sprintf("%s \n\nberikan rekomendasi menu dibawah ini untuk menjawab pertanyaan diatas\n%s \n\njawab seperti format dibawah ini\n'saya merekomendasikan bakso karena memiliki rasa yang nikmat dan gurih dan juga mengenyangkan, lalu saya merekomendasikan sop buah karena mengandung banyak buah yang menyegarkan'",
			request.Message,
			request.MenuName),
	}

	chatRequest := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{chatMessage},
	}

	resp, err := repository.ai.CreateChatCompletion(ctx, chatRequest)
	if err != nil {
		return "", err
	}

	reply := resp.Choices[0].Message.Content
	return reply, nil
}
