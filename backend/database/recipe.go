package database

import (
	"bytes"
	"encoding/json"
	"food-recipe/models"
	"net/http"
	"os"
)

var end_point = os.Getenv("HASURA_GRAPHQL_ENDPOINT")
var admin_secret = os.Getenv("HASURA_ADMIN_SECRET")
// GraphQL mutation
const insertRecipeMutation = `
	mutation ($user_id: String!, $title: String!, $description: String!, $images: [recipe_images_insert_input!]!) {
		insert_recipes_one(object: {
			user_id: $user_id,
			title: $title,
			description: $description,
			recipe_images: { data: $images }
		}) {
			id
			title
			recipe_images {
				image_url
				is_featured
			}
		}
	}
`

// InsertRecipeWithImages inserts recipe and images in Hasura
func InsertRecipeWithImages(userID string, title string, description string, images []models.RecipeImage) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"query": insertRecipeMutation,
		"variables": map[string]interface{}{
			"user_id":     userID,
			"title":       title,
			"description": description,
			"images":      images,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST",end_point, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", admin_secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if errors, ok := result["errors"]; ok {
		return nil, errors.(error)
	}

	return result["data"].(map[string]interface{}), nil
}

const updateFeaturedImageMutation = `
	mutation ($recipe_id: Int!, $new_featured_id: Int!) {
		update_recipe_images(where: {recipe_id: {_eq: $recipe_id}}, _set: {is_featured: false}) {
			affected_rows
		}
		update_recipe_images_by_pk(pk_columns: {id: $new_featured_id}, _set: {is_featured: true}) {
			id
		}
	}
`

func UpdateFeaturedImage(recipeID int, newFeaturedID int) error {
	requestBody := map[string]interface{}{
		"query": updateFeaturedImageMutation,
		"variables": map[string]interface{}{
			"recipe_id":       recipeID,
			"new_featured_id": newFeaturedID,
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", end_point, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", admin_secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}


const deleteImageMutation = `
	mutation ($image_id: Int!) {
		delete_recipe_images_by_pk(id: $image_id) {
			id
		}
	}
`

func DeleteRecipeImage(imageID int) error {
	requestBody := map[string]interface{}{
		"query": deleteImageMutation,
		"variables": map[string]interface{}{
			"image_id": imageID,
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", end_point, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", admin_secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
