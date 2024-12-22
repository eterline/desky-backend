package app

import (
	"net/http"
	"time"

	"github.com/eterline/desky-backend/internal/api"
	"github.com/eterline/desky-backend/pkg/logger"
)

func NewApp() *App {
	return &App{
		Server: api.New(),
		Log:    logger.ReturnEntry().Logger,
	}
}

func (app *App) Start() {

	var err error

	go func() {
		err = app.Server.Run()
	}()

	time.Sleep(10 * time.Minute)
	app.Server.Stop()

	if err == http.ErrServerClosed {
		app.Log.Info("closing server")
		return
	}
	app.Log.Fatalf("fatal app error: %s", err.Error())
}
