package controller

import (
	"food-recipe/models"
	"food-recipe/payments"
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
)

func InitiatePayment(c *gin.Context) {
    var request models.PaymentRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    txRef := "txn_" + request.Email

    checkoutURL, err := payments.InitiatePayment(request.Amount, request.Email, request.FirstName, request.LastName, txRef)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"checkout_url": checkoutURL})
}

// Verifies the payment status after redirect
func VerifyPayment(c *gin.Context) {
    txRef := c.Param("tx_ref")

    success, err := payments.VerifyPayment(txRef)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if success {
        c.JSON(http.StatusOK, gin.H{"message": "Payment verified successfully"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Payment not verified"})
    }
}

// Handles payment status updates (webhook)
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
