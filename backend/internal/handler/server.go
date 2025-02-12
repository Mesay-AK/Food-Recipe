package handler

import (
	// "github.com/graphql-go/handler"
	"log"
	"net/http"
	"your_project/graphql"
)

func StartGraphQLServer() {
	h := handler.New(&handler.Config{
		Schema: &graphql.Schema,
		Pretty: true,
	})

	// Setup the GraphQL endpoint
	http.Handle("/graphql", h)

	// Start the server
	log.Println("Server started at http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
