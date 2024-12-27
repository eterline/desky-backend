package handlers

import (
	"fmt"
	"net/http"

	"github.com/eterline/desky-backend/internal/configuration"
)

type APIHandlerFunc func(http.ResponseWriter, *http.Request) (op string, err error)

type BasicHandlerGroup struct {
	Config *configuration.Configuration
}

type APIResponse struct {
	StatusCode int `json:"code"`
	Message    any `json:"message"`
}

type APIErrorResponse struct {
	APIResponse
}

func (ar *APIErrorResponse) Error() string {
	return fmt.Sprintf("api error code: %d. error: %v", ar.StatusCode, ar.Message)
}

type DataErrors map[string]error
