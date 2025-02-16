package payments

import (
	"fmt"
	"food-recipe/database"
	"food-recipe/models"
	"food-recipe/handler"
)


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

	err = sendTransactionSuccessEmail(user)
	if err != nil {
		return fmt.Errorf("failed to send success email: %w", err)
	}

	return nil
}


func sendTransactionSuccessEmail(user models.User) error {
	subject := "Payment Successful"
	body := fmt.Sprintf("Dear %s,\n\nYour payment was successful. Thank you for your purchase!\n\nBest Regards,\nThe Team", user.Name)


	err := handler.SendEmail(user.Email, subject, body)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
