package authentication

import (
	"fmt"
	"time"
	"food-recipe/database"
	"food-recipe/handler"


)



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
