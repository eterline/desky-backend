package api

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/configuration"
)

type Server struct {
	HttpServer *http.Server
	Config     configuration.ServerConfig

	_ struct{}
}
