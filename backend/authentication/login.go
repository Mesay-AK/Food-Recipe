package authentication

import (
	"fmt"
	"time"
	"net/http"
	"food-recipe/database"
	"food-recipe/handler"

)



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


func Logout(refreshToken string)(error){

	err := database.DeleteRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	return nil
}