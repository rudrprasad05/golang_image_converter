package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key used for signing JWTs
var jwtKey = []byte("your_secret_key")

// User struct for storing login details
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct for JWT payload
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// In-memory user data (replace with a database in production)
var users = map[string]string{
	"admin": "password", // username: password
}

// GenerateToken generates a JWT token for a given username
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Login handler for authenticating users
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[user.Username]
	if !ok || expectedPassword != user.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Register handler for creating new users
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, exists := users[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Store the user in the in-memory map (passwords should be hashed in a real application)
	users[user.Username] = user.Password
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

// Middleware to validate JWT tokens
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// Protected handler for testing authenticated access
func Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the protected route!")
}

func convertFile(){}

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

	img, decodeErr := jpeg.Decode(file)
	if decodeErr != nil {
		http.Error(w, "Error decoding JPG image", http.StatusBadRequest)
		return
	}

	src, _ := UploadFileToS3(file, handler.Filename, "mctechfiji")

	log.Println(src)
	log.Println(img)


	// tempFile, err := os.CreateTemp("uploads", "upload-*.png")
	// if err != nil {
	// 	http.Error(w, "Could not create temporary file", http.StatusInternalServerError)
	// 	return
	// }
	// defer tempFile.Close()

	// if err := png.Encode(tempFile, img); err != nil {
	// 	http.Error(w, "Error encoding image to PNG", http.StatusInternalServerError)
	// 	return
	// }
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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", Register)
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/protected", Authenticate(Protected))
	mux.HandleFunc("/upload", UploadFileHandler)
	mux.HandleFunc("/download", DownloadImageHandler)

	handler := enableCORS(mux)

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
