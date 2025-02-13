package authentication

import (
	"fmt"
	"time"
	"net/http"
	"food-recipe/internal/database"
	"food-recipe/internal/handler"
	"food-recipe/internal/models"

)

// RegisterUser - Registers a new user
func RegisterUser(name, email, password string) (string, string, error) {
	user := models.User{
		Name: name,
		Email:    email,
		Password: password,
	}

	// Validate the user data
	err := models.ValidateUser(user)
	if err != nil {
		return "", "", err
	}

	// Check if the user already exists
	_, err = database.GetUserByEmail(email)
	if err == nil {
		return "", "", fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := handler.HashPassword(password)
	if err != nil {
		return "", "", err
	}
	user.Password = hashedPassword

	// Insert user into the database
	_, err = database.InsertUserIntoHasura(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to insert user: %w", err)
	}

	// Generate access and refresh tokens
	accessToken, err := handler.GenerateJWT(24*time.Hour, false, user.Name, user.ID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := handler.GenerateJWT(7*24*time.Hour, true, user.Name, user.ID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store the refresh token in the database
	err = database.StoreRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}


func LoginUser(email, password string, w http.ResponseWriter) (string, string, error) {
	user, err := database.GetUserByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}


	err = handler.VerifyPassword(password, user.Password)
	if err != nil {
		return "", "", fmt.Errorf("incorrect password")
	}

	accessToken, err := handler.GenerateJWT(24*time.Hour, false, user.Name, user.ID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := handler.GenerateJWT(7*24*time.Hour, true, user.Name, user.ID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}


	err = database.StoreRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	userID, name, email, err := database.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	newToken, err := handler.GenerateJWT(24*time.Hour, false, name, userID, email)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newToken, nil
}

func ForgotPassword(email models.Email)(error){

	user, err := database.GetUserByEmail(email.Email)
	if err != nil {
		return err
	}

	resetToken, err := handler.GenerateResetToken()
	if err != nil {
		return err
	}

	err = database.StorePasswordResetRequest(user.ID, resetToken); 
	if err != nil {
		return err
	}

	err = handler.SendResetEmail(user.Email, resetToken)
	if err != nil {
		return err
	}
	return nil
}


func ResetPassword(requestData models.ResetrequestData)(error){


	userID, err := database.ValidatePasswordResetToken(requestData.Token)
	if err != nil {
		return err
	}

	hashedPassword, err := handler.HashPassword(requestData.Password)
	if err != nil {
		return err
	}

	if err := database.UpdateUserPassword(userID, hashedPassword); err != nil {
		return err
	}

	return nil
}

func Logout(refreshToken string)(error){

	err := database.DeleteRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	return nil
}