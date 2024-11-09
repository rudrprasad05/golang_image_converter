package main

import (
	"backend/lib"
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

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

func GetFileFromS3(key string, format string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return "error fetching", err
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Fetch the image from S3
	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("mctechfiji"),
		Key:    aws.String(key),
	})
	if err != nil {
		return "error fetching", err
	}
	defer resp.Body.Close()

	// Decode the image using the imaging package
	srcImage, err := jpeg.Decode(resp.Body)
	if err != nil {
		return "error decoding img", err
	}

	// Create a buffer to store the converted image
	var buf bytes.Buffer

	// Convert the image to the desired format
	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(&buf, srcImage)
		if err != nil {
			return "unsuporrted imgage format", err
		}
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, srcImage, nil)
		if err != nil {
			return "unsuporrted imgage format", err
		}
	default:
		return "unsuporrted imgage format", err
	}

	// Upload the converted image back to S3
	convertedKey := strings.TrimSuffix(key, ".jpg") + "_converted." + format
	_, uploadErr := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("mctechfiji"),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("image/" + format),
	})

	if uploadErr != nil {
		log.Printf("Error uploading converted image to S3: %v", uploadErr)
		return "upload err", uploadErr
	}

	// Return the URL of the converted image
	convertedURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "mctechfiji", convertedKey)

	return convertedURL, nil
}

func DownloadImageFromS3(key string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Fetch the image from S3
	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("mctechfiji"),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading S3 object body: %w", err)
	}

	return buf.Bytes(), nil
}

func UploadFileToS3(file io.Reader, handler *multipart.FileHeader) (string, error) {
	fileName := "converted_" + lib.GenerateRandomName() + filepath.Ext(handler.Filename)
	url := "image_converter/" + fileName

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("error: %v", err)
		return "unable to read config", err
	}

	contentType := handler.Header.Get("Content-Type")
	if contentType == "" {
		// Fallback: Use the file extension to determine the MIME type
		extension := filepath.Ext(fileName)
		contentType = mime.TypeByExtension(extension)
		if contentType == "" {
			contentType = "application/octet-stream" // Default content type
		}
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	_, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("mctechfiji"),
		Key:    aws.String(url),
		Body:   file,
		ContentType: aws.String(contentType),
	})

	if uploadErr != nil{
		log.Printf("error2: %v", uploadErr)
		return "couldnt upload", uploadErr
	}

	fullUrl := "https://mctechfiji.s3.us-east-1.amazonaws.com/" + url

	return fullUrl, nil
}
