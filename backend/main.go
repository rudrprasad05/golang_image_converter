package main

import (
	"backend/lib"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"strings"
)

// Protected handler for testing authenticated access
func Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the protected route!")
}

func ConvertFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB limit

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	defer file.Close() 

	log.Println("File:", file, "Handler:", handler.Filename)

	fileType := r.FormValue("type")
	log.Println("Type:", fileType)

	srcImage, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Unsupported image format", http.StatusBadRequest)
		return
	}

	// Buffer to store the encoded image
	var buf bytes.Buffer

	err = lib.EncodeImage(&buf, srcImage, fileType)
	if err != nil {
		http.Error(w, "Error encoding image", http.StatusInternalServerError)
		return
	}

	handler, metaDataErr := lib.GetFileMetadata(fileType, &buf)
	if metaDataErr != nil {
		http.Error(w, "Error setting image headers", http.StatusInternalServerError)
		return
	}

	convertedSrc, uploadErr := UploadFileToS3(&buf, handler)

	if uploadErr != nil {
		http.Error(w, "Error uploading image", http.StatusBadRequest)
		return
	}

	log.Println(handler.Header)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(convertedSrc)

}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data with a maximum file size limit (e.g., 10MB)
	r.ParseMultipartForm(10 << 20) // 10MB limit

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "An internal server error occurred", http.StatusBadRequest)
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
	defer file.Close()

	_, seekErr := file.Seek(0, io.SeekStart)
	if seekErr != nil {
		http.Error(w, "Error resetting file pointer", http.StatusInternalServerError)
		return
	}

	src, uploadErr := UploadFileToS3(file, handler)

	if uploadErr != nil {
		http.Error(w, "Error uploading image", http.StatusBadRequest)
		return
	}

	// Send the URL back to the user as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, src) // Return the URL in JSON format

	log.Println("File uploaded successfully:", src)
}

func DownloadImageHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	fileNameParts := strings.Split(fileName, "image_converter")
	fileName = "image_converter" + fileNameParts[len(fileNameParts) - 1]
	log.Println(fileName)
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	file, err := DownloadImageFromS3(fileName)
	if err != nil{
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}

	// Set headers to prompt the download
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(file)
	if writeErr != nil {
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", UploadFileHandler)
	mux.HandleFunc("/convert", ConvertFile)
	mux.HandleFunc("/download", DownloadImageHandler)

	handler := lib.EnableCORS(mux)

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
