package models

import (
	"fmt"
	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required"`
}

type LoginDetails struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResetrequestData struct {
	Token    string `json:"token" validate:"required,min=5"`
	Password string `json:"password" validate:"required,min=8"`
}


type Email struct {
	Email string `json:"email" validate:"required,email"`
}


func ValidateUser(user User) error {
	err := validate.Struct(user)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			firstErr := validationErrors[0]
			return fmt.Errorf("field '%s' failed validation: %s", firstErr.Field(), firstErr.Tag())
		}
		return err
	}
	return nil
}
