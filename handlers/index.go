package groupie

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

// Define a struct to match the structure of the API response
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Members      []string `json:"members"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// fetchArtistData makes an HTTP GET request to the API and retrieves artist data.
func FetchArtistData() ([]Artist, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// IndexHandler handles the main page rendering and calls fetchArtistData for data retrieval.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	artists, err := FetchArtistData()
	if err != nil {
		log.Printf("Failed to fetch artist data: %s", err)
		ErrorHandler(w, r, http.StatusInternalServerError, []string{"Internal Server Error"})
		return
	}

	// Load and parse the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Failed to open template index.html: %s", err)
		ErrorHandler(w, r, http.StatusInternalServerError, []string{"Internal Server Error"})
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, artists)
	if err != nil {
		log.Printf("Failed to execute template: %s", err)
		ErrorHandler(w, r, http.StatusInternalServerError, []string{"Internal Server Error"})
		return
	}
}
