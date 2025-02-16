package authentication

import (
	"fmt"
	"food-recipe/database"
	"food-recipe/handler"
	"food-recipe/models"

)


func RegisterUser(name, email, password string) (error) {

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "user",  
	}

	err := models.ValidateUser(user)
	if err != nil {
		return err
	}

	_, err = database.GetUserByEmail(email)
	if err == nil {
		return fmt.Errorf("email already exists")
	}

	hashedPassword, err := handler.HashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	newUser, err := database.InsertUserIntoHasura(user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	subject := "Welcome to Food Recipe App!"
	body := fmt.Sprintf(
		"Hi %s,\n\nWelcome to Food Recipe App! We're excited to have you on board.\n\nHappy Cooking!\n\nBest regards,\nFood Recipe Team",
		newUser.Name,
	)

	if err := handler.SendEmail(newUser.Email, subject, body); err != nil {
		return fmt.Errorf("Failed to send welcome email: %v", err)
	}

	return nil
}



