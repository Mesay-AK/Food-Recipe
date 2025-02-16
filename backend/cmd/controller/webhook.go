package controller

import (
	"log"
	"net/http"
	"food-recipe/payments"
	"github.com/gin-gonic/gin"
)


func Webhooks(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		log.Println("Error reading Chapa webhook body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	err := payments.ProcessPayment(webhookData)
	if err != nil {
		log.Println("Error processing payment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
