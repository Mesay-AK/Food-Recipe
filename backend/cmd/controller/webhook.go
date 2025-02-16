package controller

import (
    "github.com/gin-gonic/gin"
    "net/http"
)


func HandleWebhook(c *gin.Context) {
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
