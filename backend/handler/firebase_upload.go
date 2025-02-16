package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"
	// "cloud.google.com/go/storage"
	"food-recipe/config"
)

const BucketName = "your-firebase-bucket-name"

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader, userID string) (string, error) {
	ctx := context.Background()
	app := config.FirebaseApp

	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get Firebase storage client: %w", err)
	}

	bucket, err := client.Bucket(BucketName)
	if err != nil {
		return "", fmt.Errorf("failed to get Firebase bucket: %w", err)
	}

	fileName := fmt.Sprintf("recipes/%s/%d_%s", userID, time.Now().Unix(), fileHeader.Filename)
	wc := bucket.Object(fileName).NewWriter(ctx)
	wc.ContentType = fileHeader.Header.Get("Content-Type")

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	if _, err := wc.Write(buf.Bytes()); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", BucketName, fileName)
	return url, nil
}


func DeleteImage(imagePath string) error {
	ctx := context.Background()
	app := config.FirebaseApp

	client, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Firebase storage client: %w", err)
	}

	bucket, err := client.Bucket(BucketName)
	if err != nil {
		return fmt.Errorf("failed to get Firebase bucket: %w", err)
	}

	object := bucket.Object(imagePath)
	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}