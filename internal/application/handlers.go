package application

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/services/router/handler"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	handler.StatusOK(w, "health")
}

// ====================================================================
