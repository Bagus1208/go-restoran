package helper

import (
	"restoran/features/admin/model"
)

func RequestToAdmin(data model.AdminInput) *model.Admin {
	return &model.Admin{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
}

func AdminToResponse(data *model.Admin) *model.AdminResponse {
	return &model.AdminResponse{
		ID:       data.ID,
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
}
