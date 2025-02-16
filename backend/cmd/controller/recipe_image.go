package controller

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"food-recipe/images"
	"strconv"
)

func UploadRecipeImage(c *gin.Context) {
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


func UpdateFeaturedImage(c *gin.Context) {
	recipeID, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	newFeaturedID, err := strconv.Atoi(c.Param("new_featured_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	err = images.UpdateFeaturedImage(recipeID, newFeaturedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update featured image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Featured image updated successfully"})
}


func DeleteRecipeImage(c *gin.Context) {
	// userID := c.MustGet("user_id").(string) 

	imageID, err := strconv.Atoi(c.Param("image_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	imageURL := c.Query("image_url") // Pass the image URL as a query param

	err = images.DeleteRecipeImage(imageID, imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
