package model

type AdminResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCredential struct {
	Name   string         `json:"name"`
	Access map[string]any `json:"access"`
}
