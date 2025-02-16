package authentication

import (
	"fmt"
	"food-recipe/database"
	"food-recipe/handler"
	"food-recipe/models"

)
func RegisterUser(name, email, password string) (string, string, error) {

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "user",  
	}

	err := models.ValidateUser(user)
	if err != nil {
		return "", "", err
	}

	_, err = database.GetUserByEmail(email)
	if err == nil {
		return "", "", fmt.Errorf("email already exists")
	}

	hashedPassword, err := handler.HashPassword(password)
	if err != nil {
		return "", "", err
	}
	user.Password = hashedPassword

	_, err = database.InsertUserIntoHasura(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to insert user: %w", err)
	}

	return "Registration successful! Please log in.", "", nil
}



