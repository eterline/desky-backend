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
	log *logrus.Logger = nil
)

func NewServer() *Server {
	config := configuration.GetConfig().Server
	routing := routing.InitAPIRouting(1)
	cache.Init()

	server := &Server{
		HttpServer: &http.Server{
			Addr:    config.Address(),
			Handler: routing.ConfigRoutes(),
		},
	}

	return server
}

func (srv *Server) Run() error {
	log = logger.ReturnEntry().Logger

	conf := configuration.GetConfig()

	log.Infof("required auth status: %v", conf.Auth.Enabled)
	log.Infof("starting server at: %s", conf.Server.Address())
	log.Infof("app page at: %s", conf.Server.PageAddr())

	return func(c configuration.ServerTLSConfig, s *http.Server) error {

		if c.Enabled {
			log.Infof("tls enabled. key file: %s cert file %s", c.Key, c.Certificate)
			return s.ListenAndServeTLS(c.Certificate, c.Key)
		}
		return s.ListenAndServe()

	}(conf.Server.TLS, srv.HttpServer)
}

func (srv *Server) Stop() error {
	return srv.HttpServer.Shutdown(context.Background())
}
