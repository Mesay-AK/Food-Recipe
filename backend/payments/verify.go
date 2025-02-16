package payments

import (
	"encoding/json"
	"errors"
	"food-recipe/models"
	"net/http"
	"os"
)

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
