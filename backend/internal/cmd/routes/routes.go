package routes

import (
	"food-recipe/internal/cmd/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoutes()*gin.Engine {
    router := gin.Default()
	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)
	router.POST("/refresh-token", controller.RefreshToken)

	return router
}
