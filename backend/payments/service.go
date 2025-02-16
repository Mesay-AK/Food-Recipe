package payments

import (
	"bytes"
	"encoding/json"
	"errors"
	"food-recipe/models"
	"net/http"
	"os"
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
