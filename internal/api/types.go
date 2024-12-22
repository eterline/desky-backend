package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server struct {
	HttpServer http.Server
	log        *logrus.Logger

	_ struct{}
}
