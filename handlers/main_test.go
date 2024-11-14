package groupie

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func setupMockCacheForFilteredArtistsHandler() {
	dataCache = DataCache{
		Artists: []CachedArtist{
			{
				Artist: Artist{
					Name:         "The Test Band",
					Members:      []string{"Alice", "Bob"},
					FirstAlbum:   "1995-06-15",
					CreationDate: 1990,
				},
				Locations: []string{"New York", "Los Angeles"},
			},
			{
				Artist: Artist{
					Name:         "Sample Artist",
					Members:      []string{"Charlie", "Dana"},
					FirstAlbum:   "2001-11-23",
					CreationDate: 2000,
				},
				Locations: []string{"Chicago", "Houston"},
			},
		},
		LastFetched: time.Now(),
	}
}

func TestFilteredArtistsHandler(t *testing.T) {
	// Set up mock cache data
	setupMockCacheForFilteredArtistsHandler()

	tests := []struct {
		name         string
		query        string
		wantStatus   int
		wantResponse []CachedArtist
	}{
		{
			name:       "Query matches artist name",
			query:      "Test Band",
			wantStatus: http.StatusOK,
			wantResponse: []CachedArtist{
				{
					Artist: Artist{
						Name:         "The Test Band",
						Members:      []string{"Alice", "Bob"},
						FirstAlbum:   "1995-06-15",
						CreationDate: 1990,
					},
					Locations: []string{"New York", "Los Angeles"},
				},
			},
		},
		{
			name:       "Query matches member name",
			query:      "Charlie",
			wantStatus: http.StatusOK,
			wantResponse: []CachedArtist{
				{
					Artist: Artist{
						Name:         "Sample Artist",
						Members:      []string{"Charlie", "Dana"},
						FirstAlbum:   "2001-11-23",
						CreationDate: 2000,
					},
					Locations: []string{"Chicago", "Houston"},
				},
			},
		},
		{
			name:         "Empty query parameter",
			query:        "",
			wantStatus:   http.StatusBadRequest,
			wantResponse: nil, // Expecting no response body, just a 400 status code
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/filtered-artists?q="+tt.query, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(FilteredArtistsHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("FilteredArtistsHandler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			// If expecting a JSON response, decode and compare it
			if tt.wantStatus == http.StatusOK && tt.wantResponse != nil {
				var got []CachedArtist
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}
				if !reflect.DeepEqual(got, tt.wantResponse) {
					t.Errorf("FilteredArtistsHandler returned unexpected body: got %v want %v", got, tt.wantResponse)
				}
			} else if tt.wantStatus == http.StatusBadRequest {
				expectedErrMsg := "Search query is required\n"
				if rr.Body.String() != expectedErrMsg {
					t.Errorf("FilteredArtistsHandler returned unexpected error message: got %v want %v", rr.Body.String(), expectedErrMsg)
				}
			}
		})
	}
}

func setupMockCache() {
	dataCache = DataCache{
		Artists: []CachedArtist{
			{
				Artist: Artist{
					Name:         "Test Artist",
					Members:      []string{"Member One", "Member Two"},
					FirstAlbum:   "2000-01-01",
					CreationDate: 1990,
					Locations:    "/location",
				},
				Locations: []string{"Location A", "Location B"},
			},
		},
		LastFetched: time.Now(),
	}
}

func TestSearchHandler(t *testing.T) {
	// Set up mock cache data for testing
	setupMockCache()

	tests := []struct {
		name         string
		query        string
		wantStatus   int
		wantResponse []SearchResult
	}{
		{
			name:       "Valid query parameter - match artist name",
			query:      "Test Artist",
			wantStatus: http.StatusOK,
			wantResponse: []SearchResult{
				{Name: "Test Artist", Category: "artist/band"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with or without a query parameter
			req, err := http.NewRequest("GET", "/search?q="+tt.query, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			// Use ResponseRecorder to capture the response
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SearchHandler)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("SearchHandler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			// Check the response body if the expected response is not nil
			if tt.wantStatus == http.StatusOK && tt.wantResponse != nil {
				var got []SearchResult
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}
				if !reflect.DeepEqual(got, tt.wantResponse) {
					t.Errorf("SearchHandler returned unexpected body: got %v want %v", got, tt.wantResponse)
				}
			} else if tt.wantStatus == http.StatusBadRequest {
				expectedErrMsg := "Search query is required\n"
				if rr.Body.String() != expectedErrMsg {
					t.Errorf("SearchHandler returned unexpected error message: got %v want %v", rr.Body.String(), expectedErrMsg)
				}
			}
		})
	}
}

func TestDatesHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		mockStatusCode int
		expectedStatus int
	}{
		{
			name:           "Invalid artist ID",
			query:          "?id=abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Artist ID found",
			query:          "?id=1",
			mockResponse:   `{"index": [{"id": 1, "dates": ["2023-09-12", "2023-10-01"]}]}`,
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with appropriate query parameters
			req, err := http.NewRequest("GET", "/dates"+tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Mock response recorder
			rr := httptest.NewRecorder()

			// Create a mock server to serve the external API response
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			// Override the client temporarily
			originalClient := http.DefaultClient
			http.DefaultClient = mockServer.Client()
			defer func() { http.DefaultClient = originalClient }()

			// Call the handler
			DatesHandler(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestLocationsHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		mockStatusCode int
		expectedStatus int
	}{
		{
			name:           "Artist ID found",
			query:          "?id=1",
			mockResponse:   `{"index": [{"id": 1, "locations": ["New York", "Los Angeles"], "dates": "2023-09-12"}]}`,
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/locations"+tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			rr := httptest.NewRecorder()
			// Mock HTTP server to return the appropriate response
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()
			// Replace the external API call with a call to the mock server
			http.DefaultClient = mockServer.Client()
			LocationsHandler(rr, req)
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestRelationHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		mockStatusCode int
		expectedStatus int
	}{
		{
			name:           "Artist ID found",
			query:          "?id=1",
			mockResponse:   `{"index": [{"id": 1, "datesLocations": {"New York": ["2023-09-12"], "Los Angeles": ["2023-09-15"]}}]}`,
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request
			req, err := http.NewRequest("GET", "/relation"+tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			// Create a response recorder to capture the handler's response
			rr := httptest.NewRecorder()
			// Mock HTTP server to return the appropriate mock response
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()
			// Replace the external API call with a call to the mock server
			// To simulate calling the real API endpoint but with a mock response
			http.DefaultClient = mockServer.Client()
			// Call the handler
			RelationHandler(rr, req)
			// Check if the status code is what we expect
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
