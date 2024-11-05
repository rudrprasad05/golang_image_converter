package main

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func UploadFileToS3(file multipart.File, fileName string, bucketName string) (string, error) {
	// Load the AWS configuration

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("error: %v", err)
		return "unable to read config", err
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	url := "image_converter/" + fileName


	_, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("mctechfiji"),
		Key:    aws.String(url),
		Body:   file,
	})

	if uploadErr != nil{
		log.Printf("error2: %v", uploadErr)
		return "couldnt upload", uploadErr
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, url)

	return url, nil
}
