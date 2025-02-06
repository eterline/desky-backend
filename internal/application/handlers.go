package application

import (
	"context"
	"net/http"

	"github.com/eterline/desky-backend/internal/services/router/handler"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	handler.StatusOK(w, "healthy")
}

// Useless LOL ===================================================

func PowerOFFHandler(stop context.CancelFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer stop()
		handler.StatusOK(w, "shutting down service")
	}
}

func RebootHandler(stop context.CancelFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer stop()
		handler.StatusOK(w, "restart service")
	}
}

// ====================================================================

func PreferencesHandler(auth bool, darkTheme bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pref := PresencesResponse{
			DarkTheme:  darkTheme,
			Language:   "EN",
			Background: "none",
			Auth:       auth,
		}

		handler.WriteJSON(w, http.StatusOK, pref)
	}
}
