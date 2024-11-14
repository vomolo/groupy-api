package groupie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SearchResult defines the structure for each suggestion with category details.
type SearchResult struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// LocationData holds the structure for the location data
type LocationData struct {
	Locations []string `json:"locations"`
}

// CachedArtist includes artist data and cached location data
type CachedArtist struct {
	Artist    Artist
	Locations []string
}

// Cache structure to store both artist and location data and timestamp.
type DataCache struct {
	Artists     []CachedArtist
	LastFetched time.Time
}

// Cache duration (10 minutes)
const CacheDuration = 20 * time.Minute

// Cache variable to hold the artist and location data.
var dataCache = DataCache{}

// Preload data cache on server start
func init() {
	go func() {
		if _, err := PreloadDataCache(); err != nil {
			log.Printf("Error preloading cache: %v", err)
		}
	}()
}

// PreloadDataCache fetches artist and location data on server startup
func PreloadDataCache() ([]CachedArtist, error) {
	return FetchArtistDataWithLocations()
}

// FetchLocations fetches location data for a given URL.
func FetchLocations(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read location data response: %v", err)
	}

	var locationData LocationData
	if err := json.Unmarshal(body, &locationData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location data: %v", err)
	}

	return locationData.Locations, nil
}

// FetchArtistDataWithLocations fetches artist data along with location data and updates the cache.
func FetchArtistDataWithLocations() ([]CachedArtist, error) {
	artists, err := FetchArtistData()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artist data: %v", err)
	}

	var wg sync.WaitGroup
	var cachedArtists []CachedArtist
	cachedArtistChan := make(chan CachedArtist, len(artists))

	for _, artist := range artists {
		wg.Add(1)
		go func(artist Artist) {
			defer wg.Done()
			locations, err := FetchLocations(artist.Locations)
			if err != nil {
				log.Printf("Error fetching locations for artist %s: %v", artist.Name, err)
				locations = []string{}
			}
			cachedArtistChan <- CachedArtist{Artist: artist, Locations: locations}
		}(artist)
	}

	// Wait for all goroutines to complete and close the channel
	go func() {
		wg.Wait()
		close(cachedArtistChan)
	}()

	for cachedArtist := range cachedArtistChan {
		cachedArtists = append(cachedArtists, cachedArtist)
	}

	// Update cache
	dataCache = DataCache{
		Artists:     cachedArtists,
		LastFetched: time.Now(),
	}

	return cachedArtists, nil
}

// SearchHandler handles search functionality and returns categorized suggestions.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Check if cache is expired
	if time.Since(dataCache.LastFetched) > CacheDuration {
		// Cache expired, fetch new artist and location data
		_, err := FetchArtistDataWithLocations()
		if err != nil {
			log.Printf("Failed to fetch artist data with locations: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Use cached data
	cachedArtists := dataCache.Artists
	var suggestions []SearchResult
	query = strings.ToLower(query)

	for _, cachedArtist := range cachedArtists {
		artist := cachedArtist.Artist

		// Check artist/band name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			suggestions = append(suggestions, SearchResult{
				Name:     artist.Name,
				Category: "artist/band",
			})
		}

		// Check members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				suggestions = append(suggestions, SearchResult{
					Name:     member,
					Category: "member",
				})
			}
		}

		// Check locations
		for _, location := range cachedArtist.Locations {
			if strings.Contains(strings.ToLower(location), query) {
				suggestions = append(suggestions, SearchResult{
					Name:     location,
					Category: "location",
				})
			}
		}

		// Check first album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			suggestions = append(suggestions, SearchResult{
				Name:     artist.FirstAlbum,
				Category: "first album date",
			})
		}

		// Check creation date
		creationDateStr := strconv.Itoa(artist.CreationDate)
		if strings.Contains(creationDateStr, query) {
			suggestions = append(suggestions, SearchResult{
				Name:     creationDateStr,
				Category: "creation date",
			})
		}
	}

	// Convert suggestions to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		log.Printf("Failed to encode search suggestions: %s", err)
		http.Error(w, "Failed to return search suggestions", http.StatusInternalServerError)
	}
}
