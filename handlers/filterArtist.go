package groupie
import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)
// FilteredArtistsHandler fetches and returns all artist data matching the search query.
func FilteredArtistsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}
	// Refresh cache if expired
	if time.Since(dataCache.LastFetched) > CacheDuration {
		if _, err := FetchArtistDataWithLocations(); err != nil {
			log.Printf("Failed to refresh artist data with locations: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	query = strings.ToLower(query)
	var filteredArtists []CachedArtist
	// Filter through cached artist data based on the search query
	for _, cachedArtist := range dataCache.Artists {
		artist := cachedArtist.Artist
		matchFound := false
		// Check artist name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			matchFound = true
		}
		// Check members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				matchFound = true
				break
			}
		}
		// Check locations
		for _, location := range cachedArtist.Locations {
			if strings.Contains(strings.ToLower(location), query) {
				matchFound = true
				break
			}
		}
		// Check first album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			matchFound = true
		}
		// Check creation date
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			matchFound = true
		}
		// Add artist to the filtered list if a match was found
		if matchFound {
			filteredArtists = append(filteredArtists, cachedArtist)
		}
	}
	// Convert filtered artists to JSON and return
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredArtists); err != nil {
		log.Printf("Failed to encode filtered artists: %s", err)
		http.Error(w, "Failed to return filtered artists", http.StatusInternalServerError)
	}
}
