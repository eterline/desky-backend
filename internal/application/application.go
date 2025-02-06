package application

import (
	"context"

	"net/http/pprof"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/repository/storage"
	"github.com/eterline/desky-backend/internal/server"
	agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"
	agentclient "github.com/eterline/desky-backend/pkg/agent-client"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger = nil
)

func Exec(
	ctx context.Context,
	stopFunc context.CancelFunc,
) {
	// ================= Settings parameters =================

	log = logger.ReturnEntry().Logger
	defer log.Info("service exit")

	config := configuration.GetConfig()

	// ================= Database parameters =================

	db := storage.New(config.DB.File)

	if ok := db.Test(); !ok {
		log.Warningf("can't open db: %s. open default", config.DB.File)
	}

	if err := db.Connect(); err != nil {
		log.Panicf("db connect error: %s", err)
	}
	log.Infof("db connected to file: %s", db.Source())

	ctx = context.WithValue(ctx, "db", db)

	if err := db.MigrateTables(
		new(models.AppsTopicT),
		new(models.AppsInstancesT),
		new(models.DeskyUserT),
		new(models.ExporterInfoT),
	); err != nil {
		panic(err)
	}
	log.Infof("db migrated to: %s", db.Source())

	// ================= Services additional =================

	mon := agentmon.New(ctx)

	for _, s := range config.Services.DeskyAgent {
		log.Infof("connecting to agent: %s", s.API)
		cl, err := agentclient.Reg(s.API, s.Token)

		if err != nil {
			log.Errorf("skip: %s. error: %s", s.API, err.Error())
			continue
		}

		mon.AddSession(cl, cl.Info.Hostname, cl.Info.HostID, s.API)
		log.Infof("agent successfully connected. hostname: %s. id: %s", cl.Info.Hostname, cl.Info.HostID)
	}

	ctx = context.WithValue(ctx, "agentmon", mon)

	// ================= Server parameters =================

	srv := server.New(config.ServerSocket(), config.SSL().CertFile, config.SSL().KeyFile, config.Server.Name)

	router := server.ConfigRoutes(ctx, config)
	router.Get("/config", PreferencesHandler(false, true))
	router.Get("/health", HealthHandler)

	router.Handle("/heap", pprof.Handler("heap"))

	srv.Router(router)

	// ================= Run main server parameters =================

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
