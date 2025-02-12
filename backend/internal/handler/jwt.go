package handler

import (
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
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
	// Choose the correct secret key
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
		return "", err
	}

	return signedToken, nil
}
