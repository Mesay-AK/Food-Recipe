package routes

import (
	"food-recipe/cmd/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoutes()*gin.Engine {
    router := gin.Default()
	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)
	router.POST("/refresh-token", controller.RefreshToken)
	router.POST("/auth/forgot-password", controller.ForgotPassword)
	router.POST("/auth/reset-password", controller.ResetPassword)
	router.POST("/auth/logout", controller.Logout)

	return router
}
