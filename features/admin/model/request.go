package model

type AdminInput struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type AdminInputLogin struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type InputTable struct {
	TableNumber int `json:"table_number" form:"table_number" validate:"required,min=1"`
	NoTable     int `json:"no_table" form:"no_table"`
}
