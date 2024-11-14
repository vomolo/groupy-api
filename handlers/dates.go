package groupie

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Struct to hold the dates data
type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

func DatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var error []string
	// Get the artist ID from the query parameters
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		log.Printf("Missing artist ID: %d", http.StatusMethodNotAllowed)
		error = append(error, "Missing artist ID")
		ErrorHandler(w, r, http.StatusMethodNotAllowed, error)
		return
	}

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	// Make the GET request to fetch dates data
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		log.Printf("Failed to fetch data: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
	defer resp.Body.Close()

	// Read and parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	var dates Dates
	err = json.Unmarshal(body, &dates)
	if err != nil {
		log.Printf("Failed to parse JSON: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Find the dates data for the requested artist ID
	var datesData struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}
	found := false
	for _, date := range dates.Index {
		id, err := strconv.Atoi(artistID)
		if err != nil {
			log.Printf("Invalid artist ID: %s", err)
			error = append(error, "Invalid artists ID")
			ErrorHandler(w, r, http.StatusBadRequest, error)
			return
		}
		if date.ID == id {
			datesData = date
			found = true
			break
		}
	}

	// If the artist ID is not found, return an error
	if !found {
		log.Printf("Artist ID not found: %d", http.StatusBadRequest)
		error = append(error, "Artist ID not found")
		ErrorHandler(w, r, http.StatusBadRequest, error)
		return
	}

	// Return the dates data as JSON
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(datesData); err != nil {
		log.Printf("Failed to encode JSON: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
}
