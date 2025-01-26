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
	flag.BoolFunc("migrate", "To migrate default parameters to configuration file.", func(s string) error {
		if err := configuration.Migrate(configuration.FileName, 0644); err != nil {
			panic(err)
		}
		fmt.Println("Migration: default parameters migrated to file successfully:", configuration.FileName)
		os.Exit(0)
		return nil
	})

	c := flag.String("config", configuration.FileName, "Set configuration file path.")

	flag.Parse()

	if err := configuration.Init(*c); err != nil {
		panic(err)
	}

	if err := logs(); err != nil {
		panic(err)
	}
}

// @title		Desky API test
// @version	1.0
// @BasePath	/api/v1
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application.Exec(ctx)
}

func logs() error {
	return logger.InitLogger(logger.WithPretty())
}
