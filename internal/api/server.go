package api

import (
	"context"
	"net/http"

	"github.com/eterline/desky-backend/internal/api/routing"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/services/cache"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func NewServer() *Server {
	log = logger.ReturnEntry().Logger

	config := configuration.GetConfig().Server
	routing := routing.InitAPIRouting(1)

	cache.InitCache()

	server := &Server{
		HttpServer: http.Server{
			Addr:    config.Address(),
			Handler: routing.ConfigRoutes(),
		},
	}

	return server
}

func (srv *Server) Run() error {
	conf := configuration.GetConfig().Server
	server := srv.HttpServer

	log.Infof("starting server at: %s.", conf.Address())

	return func() error {
		if !conf.TLS.Enabled {
			return server.ListenAndServe()
		}

		log.Infof("tls enabled. key file: %s cert file %s", conf.TLS.Key, conf.TLS.Certificate)

		return server.ListenAndServeTLS(
			conf.TLS.Certificate,
			conf.TLS.Key,
		)
	}()
}

func (srv *Server) Stop() error {
	return srv.HttpServer.Shutdown(context.Background())
}
