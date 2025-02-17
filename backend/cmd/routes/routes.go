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
	router.POST("/inititate",controller.InitiatePayment )
	router.POST("/verify/:tx_ref",controller.VerifyPayment)
	router.POST("/webhooks", controller.Webhooks)
	router.POST("/recipes", controller.UploadRecipeImage)
	router.DELETE("/recipe-images/:image_id", controller.DeleteRecipeImage)
	router.PUT("/recipe-images/:recipe_id/featured/:new_featured_id", controller.UpdateFeaturedImage)
	router.POST("/webhook/user-registered",  controller.UserRegisteredWebhook)
	router.POST("/webhook/purchase-confirmed",  controller.PurchaseConfirmedWebhook)



	return router
}
