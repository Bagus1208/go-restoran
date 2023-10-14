package helper

import (
	"restoran/features/admin"
	"restoran/features/admin/model"
)

func RequestToAdmin(data model.AdminInput) *admin.Admin {
	return &admin.Admin{
		ID:       data.ID,
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
}
