package groupie

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Relations struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

func RelationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var error []string

	// Get the artist ID from the query parameters
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		log.Printf("Missing artist ID: %d", http.StatusBadRequest)
		error = append(error, "Missing artist ID")
		ErrorHandler(w, r, http.StatusBadRequest, error)
		return
	}

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 20 * time.Second, // 20-second timeout
	}

	// Make the GET request to fetch relation data
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Printf("Failed to fetch data: %s",err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
	defer resp.Body.Close()

	// Read and parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	var relations Relations
	err = json.Unmarshal(body, &relations)
	if err != nil {
		log.Printf("Failed to parse JSON: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Find the relation data for the requested artist ID
	var relationData struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}
	found := false
	for _, rel := range relations.Index {
		id, err := strconv.Atoi(artistID)
		if err != nil {
			log.Printf("Invalid artist ID: %s", err)
			error = append(error, "Invalid artist ID")
			ErrorHandler(w, r, http.StatusBadRequest, error)
			return
		}
		if rel.ID == id {
			relationData = rel
			found = true
			break
		}
	}

	if !found {
		log.Printf("Artist ID not found: %d", http.StatusBadRequest)
		error = append(error, "Artist ID not found")
		ErrorHandler(w, r, http.StatusBadRequest, error)
		return
	}

	// Return the relation data as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationData); err != nil {
		log.Printf("Failed to encode JSON: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
}
