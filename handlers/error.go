package groupie
import (
	"html/template"
	"log"
	"net/http"
)
type ErrorData struct {
	Code   int
	Errors []string
}
func ErrorHandler(w http.ResponseWriter, r *http.Request, code int, errors []string) {
	// Attempt to parse the template
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// Log the error and write a generic message if the template fails
		log.Printf("Failed to parse error template: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Set the response status code only if the template parsing succeeds
	w.WriteHeader(code)
	// Create the error data for the template
	data := ErrorData{Code: code, Errors: errors}
	// Attempt to execute the template
	if err := tmpl.Execute(w, data); err != nil {
		// Log the execution error
		log.Printf("Failed to execute error template: %s", err)
	}
}
