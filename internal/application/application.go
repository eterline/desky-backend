package application

import (
	"context"
	"os"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/server"
	"github.com/eterline/desky-backend/internal/services/cache"
	"github.com/sirupsen/logrus"
)

func Exec(
	ctx context.Context,
	log *logrus.Logger,
	config *configuration.Configuration,
) {
	// ================= App additional =================

	settings := new(ApplicationSettings)
	settings.SetLanguage(LangEN)
	settings.SetBG("none")
	cache.Init()

	cache.GetEntry().PushValue("bg", FileBG())

	// ================= Server parameters =================

	srv := server.New(
		config.ServerSocket(),
		config.SSL().CertFile,
		config.SSL().KeyFile,
		config.Server.Name,
	)
	router := server.ConfigRoutes(ctx, config)

	router.Get("/config", settings.SettingHandler)
	router.Get("/health", HealthHandler)
	router.Get("/api/theme", settings.ThemeHandler)
	router.Get("/api/background", settings.WriteBG)

	srv.Router(router)

	// ================= Run main server parameters =================

	go func() {
		log.Infof("server start at: %s", config.URLString())
		defer log.Info("server closed")

		if err := srv.Run(config.SSL().TLS); err != nil {
			log.Errorf("server running error: %s", err)
		}

		os.Exit(0)
	}()

	<-ctx.Done()

	if err := srv.Stop(); err != nil {
		log.Errorf("server close error: %v", err)
	}
}
