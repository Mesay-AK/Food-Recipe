package handler

import (
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"crypto/rand"
	"encoding/base64"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))
var refreshSecretKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))

type Claims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(expiry time.Duration, useRefreshSecret bool, name, userID, email string) (string, error) {

	var secretType []byte
	if useRefreshSecret {
		secretType = refreshSecretKey
	} else {
		secretType = secretKey
	}

	expirationTime := time.Now().Add(expiry)

	claims := &Claims{
		UserID: userID,
		Name:   name,
		Email:  email,
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