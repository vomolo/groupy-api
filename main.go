package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	handlers "groupie/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Use the handler function for routing
	http.HandleFunc("/", handler)
	port := ":8080"
	log.Printf("Server started on http://localhost%s", port)

	err := http.ListenAndServe(port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		handlers.IndexHandler(w, r)
	case "/locations":
		handlers.LocationsHandler(w, r)
	case "/dates":
		handlers.DatesHandler(w, r)
	case "/relations":
		handlers.RelationHandler(w, r)
	case "/search":
		handlers.SearchHandler(w, r)
	case "/getArtists":
		handlers.FilteredArtistsHandler(w, r)
	default:
		handlers.ErrorHandler(w, r, http.StatusNotFound, []string{"Page not found"})
	}
}
