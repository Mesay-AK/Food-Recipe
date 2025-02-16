package images

import (
	"errors"
	"food-recipe/database"
	"food-recipe/handler"
	"food-recipe/models"
	"mime/multipart"
	"strings"
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
			IsFeatured: i == 0, 
		})
	}


	recipeData, err := database.InsertRecipeWithImages(userID, title, description, imageURLs)
	if err != nil {
		return nil, err
	}

	return recipeData, nil
}


func UpdateFeaturedImage(recipeID int, newFeaturedID int) error {
	err := database.UpdateFeaturedImage(recipeID, newFeaturedID)
	if err != nil {
		return err
	}

	return nil
}



func DeleteRecipeImage(imageID int, imageURL string) error {

	imagePath := strings.Split(imageURL, ".com/")[1]

	err := handler.DeleteImage(imagePath)
	if err != nil {
		return err
	}


	err = database.DeleteRecipeImage(imageID)
	if err != nil {
		return err
	}

	return nil
}