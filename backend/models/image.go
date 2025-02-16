package models


type RecipeImage struct {
	ImageURL   string `json:"image_url"`
	IsFeatured bool   `json:"is_featured"`
}
