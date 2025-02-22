package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/eterline/desky-backend/internal/application"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/services/cache"
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

	if name := config.DB.File; name != "" {
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

	cache.Init()

	// ================= Database parameters =================

	{
		db := storage.New(
			storage.NewStorageSQLite(ValueDBname),
			logger.InitStorageLogger(),
		)

		err := db.Connect()
		if err != nil {
			log.Panicf("db connect error: %s", err)
		}

		defer db.Close()

		ctx = context.WithValue(ctx, "sql_database", db)
	}

	application.Exec(ctx, stop, log, config)
}
