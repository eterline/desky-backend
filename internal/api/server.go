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

	cache.Init()

	server := &Server{
		HttpServer: http.Server{
			Addr:    config.Address(),
			Handler: routing.ConfigRoutes(),
		},
	}

	return server
}

func (srv *Server) Run() error {
	conf := configuration.GetConfig()
	server := srv.HttpServer

	log.Infof("required auth status: %v", conf.Auth.Enabled)

	log.Infof("starting server at: %s", conf.Server.Address())

	return func() error {
		if !conf.Server.TLS.Enabled {
			return server.ListenAndServe()
		}

		log.Infof("tls enabled. key file: %s cert file %s", conf.Server.TLS.Key, conf.Server.TLS.Certificate)

		return server.ListenAndServeTLS(
			conf.Server.TLS.Certificate,
			conf.Server.TLS.Key,
		)
	}()
}

func (srv *Server) Stop() error {
	return srv.HttpServer.Shutdown(context.Background())
}
