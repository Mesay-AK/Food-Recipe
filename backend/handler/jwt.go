package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"food-recipe/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))
var refreshSecretKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))


func GenerateJWT(expiry time.Duration, useRefreshSecret bool, name, userID, email, role string, allowedRoles []string) (string, error) {

	var secretType []byte
	if useRefreshSecret {
		secretType = refreshSecretKey
	} else {
		secretType = secretKey
	}

	expirationTime := time.Now().Add(expiry)

	claims := &models.Claims{
		Name:  name,
		Email: email,
		HasuraClaims: models.HasuraClaims{
			UserID: userID,
			Role:   role,
			Roles:  allowedRoles,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "hasura-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretType)
	if err != nil {
		return "", fmt.Errorf("error signing the JWT token: %w", err)
	}

	return signedToken, nil
}


func GenerateResetToken() (string, error) {

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(bytes)

	return token, nil
}