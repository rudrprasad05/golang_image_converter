package routes

import (
	"backend/lib"
	"fmt"
	"net/http"
	"strings"
)


func DownloadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	fileNameParts := strings.Split(fileName, "image_converter")
	fileName = "image_converter" + fileNameParts[len(fileNameParts) - 1]

	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	file, err := lib.DownloadImageFromS3(fileName)
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