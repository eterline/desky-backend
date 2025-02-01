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
		log.Info("service exit")
		stopFunc()
	}()

	if ok := storage.TestFile(config.DB.File); !ok {
		config.DB.File = storage.DefaultName
	}

	db := storage.New(config.DB.File)
	log.Infof("db initialized in file: %s", config.DB.File)

	if err := db.Connect(); err != nil {
		log.Panicf("db connect error: %s", err)
	}
	log.Infof("db connected to file: %s", config.DB.File)

	ctx = context.WithValue(ctx, "database", db)

	if err := db.MigrateTables(); err != nil {
		panic(err)
	}
	log.Infof("db migrated to: %s", config.DB.File)

	srv := server.New(
		config.ServerSocket(),
		config.SSL().CertFile,
		config.SSL().KeyFile,
		config.Server.Name,
	)

	router := server.ConfigRoutes(ctx, config)
	router.Get("/config", PreferencesHandler(false, true))
	router.Get("/health", HealthHandler)

	srv.Router(router)

	go func() {
		log.Infof("server start at: %s", config.URLString())
		defer log.Info("server closed")

		if err := srv.Run(config.SSL().TLS); err != nil {
			log.Errorf("server running error: %s", err)
		}
	}()

	<-ctx.Done()

	if err := db.Close(); err != nil {
		log.Errorf("db close error: %s", err.Error())
	}
	if err := srv.Stop(); err != nil {
		log.Errorf("server close error: %s", err.Error())
	}
}

// a, err := agentclient.Reg("http://10.192.10.100:4000/api", "")
// if err != nil {
// 	panic(err)
// }

// data, ok := a.Info()
// if !ok {
// 	panic(ok)
// }

// mon := agentmon.New()
// mon.AddSession(a, data.Hostname, data.HostID, "http://10.192.10.100:4000/api")

// i := mon.Pool()

// go func() {
// 	for {
// 		select {
// 		case data := <-i:
// 			fmt.Println(data)
// 		}
// 	}

// }()

// fmt.Println(utils.PrettyString(mon.List()))

// <-ctx.Done()
// mon.StopPooling()
