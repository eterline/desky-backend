package api

import (
	"net/http"
)

type Server struct {
	HttpServer http.Server

	_ struct{}
}
