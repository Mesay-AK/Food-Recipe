// /graphql/resolver.go
package graphql

import (
    "github.com/graphql-go/graphql"
    "fmt"
    "food-recipe/internal/authentication"
)

// Resolver for user registration mutation
var mutationType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Mutation",
    Fields: graphql.Fields{
        "registerUser": &graphql.Field{
            Type: userType,
            Args: graphql.FieldConfigArgument{
                "username": &graphql.ArgumentConfig{Type: graphql.String},
                "email":    &graphql.ArgumentConfig{Type: graphql.String},
                "password": &graphql.ArgumentConfig{Type: graphql.String},
            },
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                username := params.Args["username"].(string)
                email := params.Args["email"].(string)
                password := params.Args["password"].(string)

                // Call the service to register the user
                token, err := authentication.RegisterUser(username, email, password)
                if err != nil {
                    return nil, fmt.Errorf("registration failed: %w", err)
                }

                return map[string]interface{}{
                    "token": token,
                }, nil
            },
        },
    },
})
