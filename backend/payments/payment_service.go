package payments

import (
	"bytes"
	"encoding/json"
	"errors"
	"food-recipe/models"
	"food-recipe/database"
	// "food-recipe/handler"
	"net/http"
	"os"
	"fmt"
)


func InitiatePayment(amount, email, firstName, lastName, txRef string) (string, error) {
    url := "https://api.chapa.co/v1/transaction/initialize"

    payload := models.PaymentRequest{
        Amount:      amount,
        Currency:    "ETB",
        Email:       email,
        FirstName:   firstName,
        LastName:    lastName,
        TxRef:       txRef,
        CallbackURL: os.Getenv("CHAPA_CALLBACK_URL"),
        ReturnURL:   os.Getenv("CHAPA_CALLBACK_URL"),
    }

    requestBody, _ := json.Marshal(payload)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAPA_SECRET_KEY"))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var paymentResp models.PaymentResponse
    json.NewDecoder(resp.Body).Decode(&paymentResp)

    if paymentResp.Status != "success" {
        return "", errors.New(paymentResp.Message)
    }

    return paymentResp.Data.CheckoutURL, nil
}

func VerifyPayment(txRef string) (bool, error) {
    url := "https://api.chapa.co/v1/transaction/verify/" + txRef

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAPA_SECRET_KEY"))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    var verifyResp models.VerifyResponse
    json.NewDecoder(resp.Body).Decode(&verifyResp)

    if verifyResp.Status != "success" || verifyResp.Data.Status != "success" {
        return false, errors.New("payment verification failed")
    }

    return true, nil
}

func ProcessPayment(webhookData map[string]interface{}) error {
	transactionStatus, ok := webhookData["status"].(string)
	if !ok || transactionStatus != "success" {
		return fmt.Errorf("payment was not successful")
	}

	userEmail, ok := webhookData["email"].(string)
	if !ok {
		return fmt.Errorf("email not found in webhook data")
	}


	user, err := database.GetUserByEmail(userEmail)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	purchase := models.Purchase{
		UserID:   user.ID,
		Amount:   webhookData["amount"].(string),
		Status:   "success",
		TxRef:    webhookData["tx_ref"].(string),
		RecipeID: 0, 
	}


	err = database.SavePurchase(purchase)
	if err != nil {
		return fmt.Errorf("failed to save purchase: %w", err)
	}


	return nil
}

