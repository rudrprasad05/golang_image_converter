package main

import (
	"backend/lib"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Protected handler for testing authenticated access
func Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the protected route!")
}

// func convertFile(){}

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

	src, uploadErr := UploadFileToS3(file, handler.Filename, "mctechfiji")

	if uploadErr != nil{
		http.Error(w, "Error uploading image", http.StatusBadRequest)
		return
	}

	log.Println(src)
}

func DownloadImageHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Construct the file path
	filePath := filepath.Join("uploads", fileName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set headers to prompt the download
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// Serve the file
	http.ServeFile(w, r, filePath)
}




func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", UploadFileHandler)
	mux.HandleFunc("/download", DownloadImageHandler)

	handler := lib.EnableCORS(mux)

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
