package application

import (
	"context"
	"net/http"
	"os"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/server"
	"github.com/eterline/desky-backend/internal/server/router/handler"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	log    *logrus.Logger               = nil
	config *configuration.Configuration = nil
)

func Exec(ctx context.Context) {
	log = logger.ReturnEntry().Logger

	config = configuration.GetConfig()

	app := Application{
		Server: server.Setup(config),
	}

	go func() {
		app.httpServe()
	}()

	<-ctx.Done()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Error(err)
	}

	log.Info("app exit")
	os.Exit(0)
}

func (app *Application) httpServe() {
	var err error

	if config.SSL().TLS {
		err = app.Server.ListenAndServeTLS(
			config.SSL().CertFile,
			config.SSL().KeyFile,
		)
	} else {
		err = app.Server.ListenAndServe()
	}

	switch err {
	case http.ErrServerClosed:
		log.Info("http server closed")
	case nil:
		log.Info("http server closed")
	default:
		log.Errorf("server running error: %s", err.Error())
	}
}

func setupRouter() http.Handler {
	defer log.Info("http server router configured")

	mux := server.ConfigRoutes()
	mux.Get("/health", healthCheckHandler)

	return mux
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	defer log.Debug("health check handler called")
	handler.StatusOK(w, "ok")
}
