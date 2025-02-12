package database

import (
    "fmt"
    "food-recipe/internal/models"
    "github.com/machinebox/graphql"
	"time"
	"context"
    "errors"
)

var client = graphql.NewClient("https://your-hasura-instance.com/v1/graphql")


var ErrUserNotFound = errors.New("user not found")

func GetUserByEmail(email string) (models.User, error) {
    req := graphql.NewRequest(`
        query($email: String!) {
            users(where: {email: {_eq: $email}}) {
                id
                username
                email
                password
            }
        }
    `)

    req.Var("email", email)

    var respData struct {
        Users []models.User `json:"users"`
    }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

    if err := client.Run(ctx, req, &respData); err != nil {
        return models.User{}, fmt.Errorf("failed to fetch user: %w", err)
    }

    if len(respData.Users) == 0 {
        return models.User{}, ErrUserNotFound
    }

    return respData.Users[0], nil
}

func InsertUserIntoHasura(user models.User) (models.User, error) {
    req := graphql.NewRequest(`
        mutation($name: String!, $email: String!, $password: String!) {
            insert_users_one(object: {name: $name, email: $email, password: $password}) {
                id
                username
                email
            }
        }
    `)

    req.Var("name", user.Name)
    req.Var("email", user.Email)
    req.Var("password", user.Password)

    var respData struct {
        InsertUser models.User `json:"insert_users_one"`
    }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
    if err := client.Run(ctx, req, &respData); err != nil {
        return models.User{}, fmt.Errorf("failed to insert user: %w", err)
    }

    return respData.InsertUser, nil
}
func StoreRefreshToken(userID, refreshToken string) error {
    expiration := time.Now().Add(time.Hour * 24 * 7)

    // Optionally, remove any previous refresh tokens for the user.
    err := RemoveOldRefreshTokens(userID)
    if err != nil {
        return fmt.Errorf("failed to remove old refresh tokens: %w", err)
    }

    req := graphql.NewRequest(`
        mutation ($user_id: uuid!, $token: String!, $expires_at: timestamptz!) {
            insert_refresh_tokens(objects: { user_id: $user_id, token: $token, expires_at: $expires_at }) {
                affected_rows
            }
        }
    `)

    req.Var("user_id", userID)
    req.Var("token", refreshToken)
    req.Var("expires_at", expiration)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.Run(ctx, req, nil); err != nil {
        return fmt.Errorf("failed to store refresh token: %w", err)
    }

    return nil
}

func RemoveOldRefreshTokens(userID string) error {
    req := graphql.NewRequest(`
        mutation ($user_id: uuid!) {
            delete_refresh_tokens(where: { user_id: { _eq: $user_id } }) {
                affected_rows
            }
        }
    `)

    req.Var("user_id", userID)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.Run(ctx, req, nil); err != nil {
        return fmt.Errorf("failed to remove old refresh tokens: %w", err)
    }

    return nil
}

func ValidateRefreshToken(token string) (string, string, string, error) {
    req := graphql.NewRequest(`
        query ($token: String!) {
            refresh_tokens(where: {token: {_eq: $token}}) {
                user {
                    id
                    name
                    email
                }
                expires_at
            }
        }
    `)

    req.Var("token", token)

    var respData struct {
        RefreshTokens []struct {
            User struct {
                ID    string `json:"id"`
                Name  string `json:"name"`
                Email string `json:"email"`
            } `json:"user"`
            ExpiresAt string `json:"expires_at"`
        } `json:"refresh_tokens"`
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.Run(ctx, req, &respData); err != nil {
        return "", "", "", fmt.Errorf("failed to validate refresh token: %w", err)
    }

    if len(respData.RefreshTokens) == 0 {
        return "", "", "", fmt.Errorf("refresh token not found")
    }

    expirationTime, err := time.Parse(time.RFC3339, respData.RefreshTokens[0].ExpiresAt)
    if err != nil || time.Now().After(expirationTime) {
        return "", "", "", fmt.Errorf("refresh token has expired")
    }

    user := respData.RefreshTokens[0].User
    return user.ID, user.Name, user.Email, nil
}
