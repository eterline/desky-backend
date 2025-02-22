package application

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/server"
	"github.com/eterline/desky-backend/internal/services/cache"
	"github.com/eterline/desky-backend/pkg/broker"
	"github.com/sirupsen/logrus"
)

func Exec(
	ctx context.Context,
	stop context.CancelFunc,
	log *logrus.Logger,
	config *configuration.Configuration,
) {
	// ================= MQTT parameters =================

	options := []broker.OptionFunc{
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
	}

	mqttBroker := broker.NewListenerWithContext(ctx, options...)

	log.Info("connecting to mqtt service")

	if err := mqttBroker.Connect(
		config.MQTTConnTimeout(),
	); err != nil {
		log.Fatalf("mqtt connection error: %v", err)
	}
	log.Info("mqtt service connected: ", config.MQTTSocket())

	// append broker to global context
	ctx = context.WithValue(ctx, "listener_mqtt", mqttBroker)

	// ================= App additional =================

	globalCache := cache.GetEntry()
	defer globalCache.EraseValues()

	globalCache.PushValue("bg", FileBG())

	settings := new(ApplicationSettings)
	settings.SetLanguage(LangEN)
	settings.SetBG("none")

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

		if ctx.Err() == nil {
			stop()
		}
	}()

	<-ctx.Done()
	// =====================================================

	log.Info("shutting down desky")

	if err := srv.Stop(); err != nil {
		log.Errorf("server close error: %v", err)
	}
}
