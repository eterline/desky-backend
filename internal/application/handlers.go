package application

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/server/router/handler"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	handler.StatusOK(w, "healthy")
}

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
