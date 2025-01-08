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
		Config: config,
	}

	return server
}

func (srv *Server) Run() error {
	log = logger.ReturnEntry().Logger
	conf := srv.Config

	log.Infof("starting server at: %s", conf.Address())
	log.Infof("app page at: %s", conf.PageAddr())

	return func(tls configuration.ServerTLSConfig, srv *http.Server) error {

		if tls.Enabled {
			log.Infof("tls enabled. key file: %s cert file %s", tls.Key, tls.Certificate)
			return srv.ListenAndServeTLS(tls.Certificate, tls.Key)
		}
		return srv.ListenAndServe()

	}(conf.TLS, srv.HttpServer)
}

func (srv *Server) Stop() error {
	return srv.HttpServer.Shutdown(context.Background())
}
