package api

import (
	"context"
	"net/http"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/api/middlewares"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-chi/chi"
)

func New() *Server {
	config := configuration.GetConfig().Server
	chi := chi.NewMux()

	ConfigureRoutes(chi)

	server := &Server{
		HttpServer: http.Server{
			Addr:    config.Address(),
			Handler: chi,
		},
		log: logger.ReturnEntry().Logger,
	}

	return server
}

func (srv *Server) Run() error {
	conf := configuration.GetConfig().Server
	srv.log.Infof("starting server at: %s.", conf.Address())

	if conf.TLS.Enabled {

		srv.log.Infof(
			"TLS enabled. Certificate file: %s. Key file: %s",
			conf.TLS.Certificate, conf.TLS.Key,
		)
		return srv.HttpServer.ListenAndServeTLS(conf.TLS.Certificate, conf.TLS.Key)
	}

	srv.log.Infof("TLS disabled")
	return srv.HttpServer.ListenAndServe()
}

func (srv *Server) Stop() error {
	return srv.HttpServer.Shutdown(context.Background())
}

func ConfigureRoutes(rt *chi.Mux) {
	log := logger.ReturnEntry()

	mw := middlewares.Init()
	rt.Use(mw.Logging)

	log.Info("registered middlewares")

	handler := handlers.Init()
	rt.Get("/", handler.Front)
	rt.Get("/static/*", handler.Static)
	rt.Get("/assets/*", handler.Assets)

	log.Info("registered handlers")
}
