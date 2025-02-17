package controller

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"food-recipe/handler"
)

type UserRegisteredPayload struct {
	Event struct {
		Data struct {
			New struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"new"`
		} `json:"data"`
	} `json:"event"`
}


func UserRegisteredWebhook(c *gin.Context) {
	var payload UserRegisteredPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Get user details
	userEmail := payload.Event.Data.New.Email
	userName := payload.Event.Data.New.Name

	// Send Welcome Email
	subject := "Welcome to Food Recipe!"
	body := fmt.Sprintf("Hello %s,\n\nWelcome to Food Recipe! We're excited to have you.\n\nBest,\nThe Food Recipe Team", userName)

	if err := handler.SendEmail(userEmail, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome email sent successfully"})
}

func PurchaseConfirmedWebhook(c *gin.Context) {
	var payload struct {
		Event struct {
			Data struct {
				New struct {
					Email        string `json:"email"`
					RecipeName   string `json:"recipe_name"`
					PurchaseDate string `json:"purchase_date"`
				} `json:"new"`
			} `json:"data"`
		} `json:"event"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Get purchase details
	userEmail := payload.Event.Data.New.Email
	recipeName := payload.Event.Data.New.RecipeName
	purchaseDate := payload.Event.Data.New.PurchaseDate

	// Send Confirmation Email
	subject := "Your Recipe Purchase Confirmation"
	body := fmt.Sprintf("Hello,\n\nThank you for purchasing '%s' on %s.\n\nEnjoy your meal!\n\nBest,\nThe Food Recipe Team", recipeName, purchaseDate)

	if err := handler.SendEmail(userEmail, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase confirmation email sent successfully"})
}
