// /graphql/schema.go
package graphql

import (
    "github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
    Name: "User",
    Fields: graphql.Fields{
        "id":       &graphql.Field{Type: graphql.String},
        "username": &graphql.Field{Type: graphql.String},
        "email":    &graphql.Field{Type: graphql.String},
    },
})

// Defining the full schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Mutation: mutationType,
})
