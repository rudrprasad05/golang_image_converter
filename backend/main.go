package main

import (
	"backend/lib"
	"backend/routes"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/test", routes.ConvertFile)
	mux.HandleFunc("/convert", routes.ConvertFile)
	mux.HandleFunc("/download", routes.DownloadImageHandler)

	handler := lib.EnableCORS(mux)

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
