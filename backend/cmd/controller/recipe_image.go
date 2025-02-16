package controller

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"food-recipe/images"
)

func UploadRecipe(c *gin.Context) {
	userID := c.MustGet("user_id").(string) // Extract from middleware
	title := c.PostForm("title")
	description := c.PostForm("description")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form submission"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images uploaded"})
		return
	}

	recipe, err := images.InsertRecipe(userID, title, description, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to insert recipe: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
}
