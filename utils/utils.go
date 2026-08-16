package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// ParseJSON → utility function to parse JSON payload from an HTTP request into a specified struct, returns an error if parsing fails
func ParseJSON(r *http.Request, payload any) error {

	if r.Body == nil {
		return http.ErrBodyNotAllowed
	}
	// Decode the JSON payload from the request body into the provided struct, returning any errors that occur during decoding
	return json.NewDecoder(r.Body).Decode(payload)

}

// WriteJSON → utility function to write a JSON response with a specified status code and payload, returns an error if encoding fails
func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Encode the provided payload as JSON and write it to the response, returning any errors that occur during encoding
	return json.NewEncoder(w).Encode(payload)

}

// WriteError → utility function to write a standardized JSON error response with a specified status code and error message, returns an error if encoding fails
func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]string{"error": message})
}
