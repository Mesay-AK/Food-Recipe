package images

import (
	"errors"
	"food-recipe/database"
	"food-recipe/handler"
	"food-recipe/models"
	"mime/multipart"
)


func InsertRecipe(userID string, title string, description string, files []*multipart.FileHeader) (map[string]interface{}, error) {
	if len(files) == 0 {
		return nil, errors.New("at least one image is required")
	}

	var imageURLs []models.RecipeImage

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		url, err := handler.UploadImage(file, fileHeader, userID)
		if err != nil {
			return nil, err
		}

		imageURLs = append(imageURLs, models.RecipeImage{
			ImageURL:   url,
			IsFeatured: i == 0, // First image is featured
		})
	}

	// Insert recipe and images into Hasura
	recipeData, err := database.InsertRecipeWithImages(userID, title, description, imageURLs)
	if err != nil {
		return nil, err
	}

	return recipeData, nil
}
