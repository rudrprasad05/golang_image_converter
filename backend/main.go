package main

import (
	"backend/lib"
	"backend/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rudrprasad05/go-logs/logs"
)

func main() {
	router := mux.NewRouter()
	logger, err := logs.NewLogger()

	if err != nil{
		log.Println("err", err)
		return
	}

	router.HandleFunc("/api/test", routes.TestApi).Methods("GET")
	router.HandleFunc("/convert", routes.ConvertFile)
	router.HandleFunc("/download", routes.DownloadImageHandler)

	handler := lib.EnableCORS(router)
	loggedHandler := logs.LoggingMiddleware(logger, handler)


	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", loggedHandler))
}
