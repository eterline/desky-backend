package app

import (
	"time"

	"github.com/eterline/desky-backend/internal/api"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Execute() {

	log = logger.ReturnEntry().Logger

	APIServer := api.NewServer()

	go func() {
		err := APIServer.Run()
		log.Fatalf("fatal app error: %s", err.Error())
	}()

	time.Sleep(10 * time.Minute)
	err := APIServer.Stop()

	log.Fatalf("fatal app error: %s", err.Error())
}
