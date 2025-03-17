package application

import (
	"context"
	"fmt"
	"time"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/server"
	"github.com/eterline/desky-backend/internal/services/cache"
	"github.com/eterline/desky-backend/pkg/broker"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/eterline/desky-backend/pkg/storage"
	"github.com/eterline/desky-backend/pkg/toolkit"
)

func Exec(root *toolkit.AppStarter, config *configuration.Configuration) {

	log := logger.ReturnEntry()

	globalCache := cache.GetEntry()
	globalCache.PushValue("bg", FileBG())

	settings := new(ApplicationSettings)
	settings.SetLanguage(LangEN)
	settings.SetBG("none")

	// ================= Server parameters =================

	router := server.ConfigRoutes(root.Context, config)

	router.Get("/config", settings.SettingHandler)
	router.Get("/health", HealthHandler)
	router.Get("/api/theme", settings.ThemeHandler)
	router.Get("/api/background", settings.WriteBG)

	srv := server.New(
		config.ServerSocket(),
		config.SSL().CertFile,
		config.SSL().KeyFile,
		config.Server.Name,
		router,
	)

	// ================= Run main server parameters =================

	go func() {
		log.Infof("server start at: %s", config.URLString())
		defer log.Info("server closed")

		if err := srv.Run(config.SSL().TLS); err != nil {
			log.Errorf("server running error: %s", err)
		}

		root.StopApp()
	}()

	root.Wait()

	log.Info("shutting down desky")
	if err := srv.Stop(); err != nil {
		log.Errorf("server close error: %v", err)
	}
}

func InitMqtt(ctx context.Context, testInterval time.Duration, config *configuration.Configuration) *broker.ListenerMQTT {

	log := logger.ReturnEntry()

	mqttBroker := broker.NewListenerWithContext(ctx,
		broker.OptionInsecureCerts(),

		broker.OptionClientIDString(config.Agent.UUID),
		broker.OptionServer(
			config.Agent.Server.Protocol,
			config.Agent.Server.Host,
			config.Agent.Server.Port,
		),

		broker.OptionCredentials(
			config.Agent.Username,
			config.Agent.Password,
		),

		broker.OptionDefaultQoS(config.Agent.DefaultQoS),
	)
	if err := mqttBroker.Connect(30 * time.Second); err != nil {
		log.Fatalf("broker connection error: %v", err)
	}

	go func(b *broker.ListenerMQTT) {

		connAttempts := 0
		ticks := time.NewTicker(testInterval)
		defer ticks.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticks.C:
				if b.Connected() {
					fmt.Println(2)
					continue
				}

				if err := b.Connect(30 * time.Second); err != nil {
					log.Errorf("mqtt connection attempt error: %v", err)

					if connAttempts > 4 {
						panic(err)
					}
					connAttempts++
				}
			}
		}
	}(mqttBroker)

	return mqttBroker
}

func InitDatabase() *storage.DB {
	db := storage.New(
		storage.NewStorageSQLite("desky.db"),
		logger.InitStorageLogger(),
	)

	if err := db.Connect(); err != nil {
		panic(err)
	}

	return db
}
