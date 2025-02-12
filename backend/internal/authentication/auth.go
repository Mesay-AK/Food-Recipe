package authentication

import (
	"fmt"
	"time"
	"food-recipe/internal/handler"
	"food-recipe/internal/models"
	"food-recipe/internal/database"
)

func RegisterUser(username, email, password string) (string, error) {
    
    user := models.User{
        Name: username,
        Email:    email,
        Password: password,
    }
    
    err := models.ValidateUser(user)
	if err != nil{
		return "", err
	}

    _, err = database.GetUserByEmail(email); 
    if err == nil {
        return "", fmt.Errorf("email already exists: %w", err)
    }

    hashedPassword, err := handler.HashPassword(password);
	if err != nil{
		return "", err
	}
    user.Password = hashedPassword;

    _, err = database.InsertUserIntoHasura(user)
    if err != nil {
        return "", fmt.Errorf("failed to insert user into database: %w", err)
    }

    return "user registered successfully", nil
}


func LoginUser(username, password string) (string,string, error) {
    
    user, err := database.GetUserByEmail(username)

    if err != nil {
        return "", "", fmt.Errorf("failed to find user: %w", err)
    }

    err = handler.VerifyPassword(password, user.Password)
    if err != nil{
        return "","", fmt.Errorf("incorrect password")
    }

    token, err := handler.GenerateJWT(24*time.Hour, false, user.Name, user.ID, user.Email) // 
    if err != nil {
        return "", "", fmt.Errorf("failed to generate token: %w", err)
    }
    refreshToken, err := handler.GenerateJWT(7*24*time.Hour, true, user.Name, user.ID, user.Email)
    if err != nil {
        return "","", fmt.Errorf("failed to refresh token: %w", err)
    }

    return token, refreshToken, nil
}
