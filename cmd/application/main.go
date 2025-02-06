package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eterline/desky-backend/internal/application"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/pkg/logger"
)

func init() {
	flag.BoolFunc("gen", "To generate configuration file.", genConfig)
	c := flag.String("config", configuration.FileName, "Set configuration file path.")
	flag.Parse()

	if err := configuration.Init(*c); err != nil {
		panic(err)
	}

	config := configuration.GetConfig()

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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	application.Exec(ctx, stop)
}

func genConfig(string) error {
	if err := configuration.Migrate(configuration.FileName, 0644); err != nil {
		panic(err)
	}
	fmt.Println("Migration: default config generated")
	os.Exit(0)
	return nil
}
