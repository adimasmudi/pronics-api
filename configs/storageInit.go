package configs

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type ClientUploader struct {
	client       *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func StorageInit(path string) *ClientUploader {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./project_secret_key.json")
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uploader = &ClientUploader{
		client:         client,
		bucketName: os.Getenv("BUCKET_NAME"),
		projectID:  os.Getenv("PROJECT_ID"),
		uploadPath: "images/"+path+"/",
	}

	return uploader

}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(file multipart.File, object string) (error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.client.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	
	return nil
}

// read file
func (c *ClientUploader) ReadFile(fileName string) (string,error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	return fmt.Sprintf("https://storage.cloud.google.com/%s/%s%s",c.bucketName, c.uploadPath, fileName), nil

}