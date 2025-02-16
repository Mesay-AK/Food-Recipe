package config

import (
	"context"
	"log"

	"firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	opt := option.WithCredentialsFile("path/to/your/firebase-service-account.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}

	FirebaseApp = app
}
