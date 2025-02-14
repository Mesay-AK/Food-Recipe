package chappa

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

var CHAPA_SECRET_KEY = os.Getenv("CHAPA_SECRET_KEY")

type PaymentRequest struct {
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	TxRef      string `json:"tx_ref"`
	Callback   string `json:"callback_url"`
	ReturnURL  string `json:"return_url"`
}


type PaymentResponse struct {
	Status string `json:"status"`
	Data   struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var request PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+CHAPA_SECRET_KEY).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&PaymentResponse{}).
		Post("https://api.chapa.co/v1/transaction/initialize")

	if err != nil {
		log.Println("Chapa API error:", err)
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	var response PaymentResponse
	json.Unmarshal(resp.Body(), &response)

	if response.Status != "success" {
		http.Error(w, "Payment initialization failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


