package database

import (
    "fmt"
    "food-recipe/models"
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
                role
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
                role
            }
        }
    `)

    req.Var("name", user.Name)
    req.Var("email", user.Email)
    req.Var("password", user.Password)
    req.Var("role", user.Role) 

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

func DeleteRefreshToken(token string) error {
	req := graphql.NewRequest(`
        mutation ($token: String!) {
            delete_refresh_tokens(where: {token: {_eq: $token}}) {
                affected_rows
            }
        }
    `)

	req.Var("token", token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Run(ctx, req, nil)
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

func StorePasswordResetRequest(userID, resetToken string) error {
	expirationTime := time.Now().Add(24 * time.Hour)
	req := graphql.NewRequest(`
		mutation($user_id: String!, $token: String!, $expires_at: timestamptz!) {
			insert_reset_requests(objects: {user_id: $user_id, token: $token, expiration_time: $expires_at}) {
				affected_rows
			}
		}
	`)
	req.Var("user_id", userID)
	req.Var("token", resetToken)
	req.Var("expires_at", expirationTime)

	if err := client.Run(context.Background(), req, nil); err != nil {
		return fmt.Errorf("failed to store reset request: %w", err)
	}

	return nil
}

func ValidatePasswordResetToken(token string) (string, error) {
	req := graphql.NewRequest(`
		query($token: String!) {
			reset_requests(where: {token: {_eq: $token}}) {
				user_id
				expiration_time
			}
		}
	`)
	req.Var("token", token)

	var respData struct {
		ResetRequests []struct {
			UserID        string    `json:"user_id"`
			ExpirationTime time.Time `json:"expiration_time"`
		}
	}

	if err := client.Run(context.Background(), req, &respData); err != nil {
		return "", fmt.Errorf("failed to validate reset token: %w", err)
	}

	if len(respData.ResetRequests) == 0 {
		return "", fmt.Errorf("invalid token")
	}


	if time.Now().After(respData.ResetRequests[0].ExpirationTime) {
		return "", fmt.Errorf("token expired")
	}

	return respData.ResetRequests[0].UserID, nil
}


func UpdateUserPassword(userID, hashedPassword string) error {
	req := graphql.NewRequest(`
		mutation($user_id: String!, $password: String!) {
			update_users(where: {id: {_eq: $user_id}}, _set: {password: $password}) {
				affected_rows
			}
		}
	`)
	req.Var("user_id", userID)
	req.Var("password", hashedPassword)

	if err := client.Run(context.Background(), req, nil); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}