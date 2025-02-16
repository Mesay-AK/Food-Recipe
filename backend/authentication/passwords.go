package authentication

import (

	"food-recipe/database"
	"food-recipe/handler"
	"food-recipe/models"

)

func ForgotPassword(email models.Email)(error){

	user, err := database.GetUserByEmail(email.Email)
	if err != nil {
		return err
	}

	resetToken, err := handler.GenerateResetToken()
	if err != nil {
		return err
	}

	err = database.StorePasswordResetRequest(user.ID, resetToken); 
	if err != nil {
		return err
	}

	err = handler.SendResetEmail(user.Email, resetToken)
	if err != nil {
		return err
	}
	return nil
}


func ResetPassword(requestData models.ResetrequestData)(error){


	userID, err := database.ValidatePasswordResetToken(requestData.Token)
	if err != nil {
		return err
	}

	hashedPassword, err := handler.HashPassword(requestData.Password)
	if err != nil {
		return err
	}

	if err := database.UpdateUserPassword(userID, hashedPassword); err != nil {
		return err
	}

	return nil
}
