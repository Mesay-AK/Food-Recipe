package controller

import (
	"food-recipe/models"
	"food-recipe/payments"
	"net/http"
	"github.com/gin-gonic/gin"
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

func StatusUpdate(c *gin.Context) {
        var payload map[string]interface{}
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
            return
        }

        txRef, _ := payload["tx_ref"].(string)
        status, _ := payload["status"].(string)

        if status == "success" {
            c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "tx_ref": txRef})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Payment failed"})
        }
    }

