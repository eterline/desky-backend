package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = nil

// Initialize custom type handler with error processing
func InitController(handle APIHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if log == nil {
			log = logger.ReturnEntry().Logger
		}

		op, err := handle(w, r)

		log.Debugf("requested controller: %s", op)

		if err != nil {

			switch e := err.(type) {

			case *APIErrorResponse:
				WriteJSON(w, e.StatusCode, e)
				break

			default:
				errDefault := InternalErrorResponse()
				WriteJSON(w, errDefault.StatusCode, errDefault)
				break
			}

			log.Errorf(
				"API Error - path: %s | controller: %s | error: %s",
				r.URL.Path, op, err.Error(),
			)
		}
	}
}

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
