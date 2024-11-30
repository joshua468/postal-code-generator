package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"fmt"
	"log"
	"strings"
	"os"
)

// Define a struct to hold the postal code data
type PostalCodes map[string]string

// Load postal codes from a CSV file
func loadPostalCodes() (PostalCodes, error) {
	postalCodes := make(PostalCodes)
	filePath := "./data/postal_codes.csv"

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %v", err)
	}

	// Load postal codes into the map
	for _, record := range records[1:] { // Skip the header row
		state := strings.TrimSpace(record[0])
		postalCode := strings.TrimSpace(record[1])
		postalCodes[state] = postalCode
	}

	return postalCodes, nil
}


func serveHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}
func getPostalCode(w http.ResponseWriter, r *http.Request) {
	// Get the state name from the URL
	state := strings.TrimPrefix(r.URL.Path, "/api/postal-code/")
	state = strings.TrimSpace(state)

	// Log the state being requested
	fmt.Println("Received request for state:", state)

	postalCodes, err := loadPostalCodes()
	if err != nil {
		// Log and send the error message
		fmt.Println("Error loading postal codes:", err)
		http.Error(w, "Error loading postal codes", http.StatusInternalServerError)
		return
	}

	// Check if the state exists in the postal codes map
	postalCode, found := postalCodes[state]
	if !found {
		// Log the missing state and send a not found error
		fmt.Println("State not found:", state)
		http.Error(w, fmt.Sprintf("State '%s' not found", state), http.StatusNotFound)
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
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/api/postal-code/", getPostalCode)

	// Start the server
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
