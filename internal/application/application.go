package application

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/repository/storage"
	"github.com/eterline/desky-backend/internal/server"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger = nil
)

func Exec(
	ctx context.Context,
	config *configuration.Configuration,
	stopFunc context.CancelFunc,
) {
	log = logger.ReturnEntry().Logger
	defer func() {
		stopFunc()
		log.Info("service exit")
	}()

	db := dbInitialize(config.DB.File)
	server := server.New(
		config.ServerSocket(),
		config.SSL().CertFile,
		config.SSL().KeyFile,
		"",
	)

	go func() {
		log.Infof("server start at: %s", config.ServerSocket())

		if err := server.Run(config.SSL().TLS); err != nil {
			log.Errorf("server running error: %s", err)
		}
		log.Info("server closed")
	}()

	<-ctx.Done()

	if err := db.Close(); err != nil {
		log.Errorf("db close error: %s", err.Error())
	}
	if err := server.Stop(); err != nil {
		log.Errorf("server close error: %s", err.Error())
	}
}

func dbInitialize(dbFile string) *storage.DB {

	if dbFile == "" {
		dbFile = storage.DefaultName
	}

	db := storage.InitStorage(dbFile)
	log.Infof("db initialized in file: %s", dbFile)

	if err := db.MigrateTables(); err != nil {
		panic(err)
	}
	log.Infof("db migrated to: %s", dbFile)

	err := db.Connect()
	if err != nil {
		panic(err)
	}
	log.Infof("db connected to file: %s", dbFile)

	return db
}
