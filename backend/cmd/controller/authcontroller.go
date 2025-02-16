package controller

import (
    "food-recipe/authentication"
    "food-recipe/models"
	"net/http"
	"github.com/gin-gonic/gin"
)



func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := authentication.RegisterUser(user.Name, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "registration Complete",
	})
}


func LoginUser(c *gin.Context) {
	var logInReq models.LoginDetails
	if err := c.ShouldBindJSON(&logInReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessToken, refreshToken, err := authentication.LoginUser(logInReq.Email, logInReq.Password, c.Writer)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set refresh token as cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 1 week
	})


	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
	})
}


func RefreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	newToken, err := authentication.RefreshAccessToken(cookie.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
	})
}


func ForgotPassword(c *gin.Context) {
	var email models.Email;
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func ResetPassword(c *gin.Context) {

	var requestData models.ResetrequestData;

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := authentication.ResetPassword(requestData); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while resetting password"})
	
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func Logout(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token found"})
		return
	}
	if err := authentication.Logout(refreshToken); err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return 
	}


	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
