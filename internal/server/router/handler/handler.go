package handler

import (
	"encoding/json"
	"net/http"
)

// Decode JSON request body to data structure
func DecodeRequest(r *http.Request, v any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return InvalidContentTypeResponse()
	}
	return json.NewDecoder(r.Body).Decode(&v)
}

// Send JSON object response
func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}
