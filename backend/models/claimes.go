package models

import(
	"github.com/dgrijalva/jwt-go"
)

type HasuraClaims struct {
	UserID string   `json:"x-hasura-user-id"`
	Role   string   `json:"x-hasura-role"`
	Roles  []string `json:"x-hasura-allowed-roles"`
}


type Claims struct {
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	HasuraClaims HasuraClaims `json:"https://hasura.io/jwt/claims"`
	jwt.StandardClaims
}