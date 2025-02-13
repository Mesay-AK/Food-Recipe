package models

import (
	validator "github.com/go-playground/validator/v10"
	"fmt"
)

var validate = validator.New()


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