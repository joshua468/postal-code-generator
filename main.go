package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Define a struct to hold the postal code data
type PostalCodes map[string]string

// Load postal codes from a JSON file
func loadPostalCodes() (PostalCodes, error) {
	var postalCodes PostalCodes
	filePath := filepath.Join("data", "postal_codes.json")
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &postalCodes)
	if err != nil {
		return nil, err
	}

	return postalCodes, nil
}

// Handler for serving the index.html
func serveHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

// Handler for searching postal codes
func getPostalCode(w http.ResponseWriter, r *http.Request) {
	// Get the state name from the URL
	state := strings.TrimPrefix(r.URL.Path, "/api/postal-code/")
	postalCodes, err := loadPostalCodes()
	if err != nil {
		http.Error(w, "Error loading postal codes", http.StatusInternalServerError)
		return
	}

	postalCode, found := postalCodes[state]
	if !found {
		http.Error(w, "State not found", http.StatusNotFound)
		return
	}

	// Return the postal code as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"state":       state,
		"postal_code": postalCode,
	})
}

func main() {
	// Serve static files (HTML, CSS)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// Define API routes
	http.HandleFunc("/", serveHomePage)                 // For the homepage
	http.HandleFunc("/api/postal-code/", getPostalCode) // API for postal codes

	// Start the server
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
