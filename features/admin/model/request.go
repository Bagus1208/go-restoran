package model

type AdminInput struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type AdminInputLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
