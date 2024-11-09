package lib

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"strings"
)

func GenerateRandomName() string {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %v", err)
	}
	return hex.EncodeToString(bytes) // Returns a 32-character hexadecimal string
}

func ExtractBucketAndKeyFromURL(s3URL string) (string, string) {
	// Example S3 URL: https://bucket-name.s3.amazonaws.com/path/to/object.jpg
	parts := strings.SplitN(s3URL, "/", 4)
	if len(parts) < 4 {
		return "", ""
	}
	bucket := strings.TrimPrefix(parts[2], "s3.")
	key := parts[3]
	return bucket, key
}

func GetFileMetadata(fileType string, buf *bytes.Buffer) (*multipart.FileHeader, error) {
	var contentType, extension string

	// Determine the content type and extension based on the file type
	switch strings.ToLower(fileType) {
	case "png":
		contentType = "image/png"
		extension = ".png"
	case "jpeg", "jpg":
		contentType = "image/jpeg"
		extension = ".jpg"
	case "webp":
		contentType = "image/webp"
		extension = ".webp"
	default:
		return nil, fmt.Errorf("unsupported image format: %s", fileType)
	}

	// Create and return a new FileHeader with dynamic metadata
	return &multipart.FileHeader{
		Filename: "converted_image" + extension,
		Size:     int64(buf.Len()),
		Header:   map[string][]string{"Content-Type": {contentType}},
	}, nil
}

func EncodeImage(buf *bytes.Buffer, srcImage image.Image, fileType string) error {
	switch strings.ToLower(fileType) {
	case "png":
		return png.Encode(buf, srcImage)
	case "jpeg", "jpg":
		return jpeg.Encode(buf, srcImage, nil)
	
	default:
		return fmt.Errorf("unsupported image format: %s", fileType)
	}
}