package handler

import (
	"fmt"
	"errors"
	"log"
	"net/smtp"
	"os"
	"github.com/joho/godotenv"
)

var FROM string
var SENDERPASS string
var HOST string
var PORT string

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	// Assign values to global variables (remove `var`)
	FROM = os.Getenv("EMAIL_ADDRESS")
	SENDERPASS = os.Getenv("EMAIL_SENDER_PASSWORD")
	HOST = "smtp.gmail.com"
	PORT = "587" // Use 587 for TLS instead of 465

	return nil
}

func SendEmail(to, subject, body string) error {
	err := loadEnv()
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", FROM, SENDERPASS, HOST)

	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", FROM, to, subject, body)

	err = smtp.SendMail(HOST+":"+PORT, auth, FROM, []string{to}, []byte(message))
	if err != nil {
		return errors.New("error while sending verification email: " + err.Error())
	}

	return nil
}

func SendResetEmail(email, token string) error {
	verificationLink := fmt.Sprintf("https://Food-Recipe/user/reset-password?token=%s", token)

	subject := "Reset Your Password"
	body := fmt.Sprintf("Click the link below to reset your password:\n\n%s", verificationLink)

	err := SendEmail(email, subject, body)
	if err != nil {
		return err
	}

	return nil
}
