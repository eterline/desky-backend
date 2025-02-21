package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/eterline/desky-backend/internal/application"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/pkg/broker"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/eterline/desky-backend/pkg/storage"
)

var (
	ValueDBname = "desky.db"
)

func init() {
	c := flag.String("config", configuration.FileName, "Set configuration file path.")
	flag.Parse()

	if err := configuration.Init(*c); err != nil {
		panic(err)
	}

	config := configuration.GetConfig()

	if name := os.Getenv("DESKY_DB_NAME"); name != "" {
		ValueDBname = name
	}

	if err := logger.InitLogger(
		logger.WithDevEnvBool(config.DevelopEnv),
		logger.WithPath("./logs"),
		logger.WithPretty(),
	); err != nil {
		panic(err)
	}
}

// @title		Desky API test
// @version	1.0
// @BasePath	/api/v1
func main() {

	log := logger.ReturnEntry().Logger
	config := configuration.GetConfig()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// ================= Database parameters =================

	db := storage.New(storage.NewStorageSQLite(ValueDBname), logger.InitStorageLogger())

	err := db.Connect()
	if err != nil {
		log.Panicf("db connect error: %s", err)
	}
	defer db.Close()

	{
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

			broker.OptionDefaultQoS(config.Agent.Server.DefaultQoS),
		)

		if err := mqttBroker.Connect(
			config.MQTTConnTimeout(),
		); err != nil {
			log.Fatalf("mqtt connection error: %v", err)
		}
		log.Infof("mqtt service connected: %s://%s:%d",
			config.Agent.Server.Protocol,
			config.Agent.Server.Host,
			config.Agent.Server.Port,
		)

		ctx = context.WithValue(ctx, "agentmon_mqtt", mqttBroker)
	}

	application.Exec(context.WithValue(ctx, "sql_database", db), log, config)
}
