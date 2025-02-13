package models

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginDetails struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	}

type ResetrequestData struct {
	Token    string `json:"token" validate:"required, min=5"`
	Password string `json:"password" validate:"required,min=8"`
	}

type Email struct {
		Email string `json:"email" validate:"required,email"`
	}