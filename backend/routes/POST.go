package routes

import (
	"backend/lib"
	"bytes"
	"encoding/json"
	"image"
	"net/http"
)

func TestApi(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("API Working")
}


func ConvertFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB limit

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	defer file.Close() 

	fileType := r.FormValue("type")

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

	convertedSrc, uploadErr := lib.UploadFileToS3(&buf, handler)

	if uploadErr != nil {
		http.Error(w, "Error uploading image", http.StatusBadRequest)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(convertedSrc)

}
