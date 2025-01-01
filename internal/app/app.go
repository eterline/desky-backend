package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/eterline/desky-backend/internal/api"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Execute() {

	log = logger.ReturnEntry().Logger

	log.Info("prepare server to start")
	APIServer := api.NewServer()

	process, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err := APIServer.Run()
		log.Errorf("fatal app error: %s", err.Error())

		stop()
	}()

	<-process.Done()

	log.Info("stopping the server...")

	if err := APIServer.Stop(); err != nil {
		log.Errorf("stopping app error: %s", err.Error())
	}

	log.Info("app shutdown. Bye...")
	os.Exit(0)
}
