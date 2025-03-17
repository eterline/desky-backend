package main

import (
	"time"

	"github.com/eterline/desky-backend/internal/application"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/cache"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/eterline/desky-backend/pkg/toolkit"
)

var (
	config *configuration.Configuration
)

func main() {

	root := toolkit.InitAppStart(func() error {

		if err := configuration.Init(configuration.FileName); err != nil {
			return err
		}
		config = configuration.GetConfig()

		if err := logger.InitLogger(
			logger.WithDevEnvBool(config.DevelopEnv),
			logger.WithPath("./logs"),
			logger.WithPretty(),
		); err != nil {
			return err
		}

		return nil
	})
	defer root.StopApp()

	// Initialize Database connection to app
	db := application.InitDatabase()
	defer db.Close()
	root.AddValue(models.DATABASE_CONTEXT_KEY, db)

	// Inititalize MQTT connection for app
	mqtt := application.InitMqtt(root.Context, 10*time.Minute, config)
	defer mqtt.Close()
	root.AddValue(models.MESSAGE_BROKER_CONTEXT_KEY, mqtt)

	// Init application cache singletone
	cache.Init()
	defer cache.EraseValues()

	application.Exec(root, config)
}
