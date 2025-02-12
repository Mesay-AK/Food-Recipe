package database

import (
    "fmt"
    "food-recipe/internal/models"
    "github.com/machinebox/graphql"
	"time"
	"context"
)

var client = graphql.NewClient("https://your-hasura-instance.com/v1/graphql")


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
        return models.User{}, fmt.Errorf("user not found")
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


