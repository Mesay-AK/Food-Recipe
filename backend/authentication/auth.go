package authentication

import (
	"fmt"
	"time"
	"net/http"
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


func LoginUser(email, password string, w http.ResponseWriter) (string, string, error) {
    user, err := database.GetUserByEmail(email)
    if err != nil {
        return "", "", fmt.Errorf("user not found")
    }

    err = handler.VerifyPassword(password, user.Password)
    if err != nil {
        return "", "", fmt.Errorf("incorrect password")
    }

    allowedRoles := []string{"user"} 
    if user.Role == "admin" {
        allowedRoles = append(allowedRoles, "admin") 
    }

    accessToken, err := handler.GenerateJWT(24*time.Hour, false, user.Name, user.ID, user.Email, user.Role, allowedRoles)
    if err != nil {
        return "", "", fmt.Errorf("failed to generate access token: %w", err)
    }

    refreshToken, err := handler.GenerateJWT(7*24*time.Hour, true, user.Name, user.ID, user.Email, user.Role, allowedRoles)
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


    user, err := database.GetUserByEmail(email) 
    if err != nil {
        return "", fmt.Errorf("error fetching user details: %w", err)
    }


    allowedRoles := []string{"user"}  
    if user.Role == "admin" {
        allowedRoles = append(allowedRoles, "admin") 
    }

    newToken, err := handler.GenerateJWT(24*time.Hour, false, name, userID, email, user.Role, allowedRoles)
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