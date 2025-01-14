package main

import (
	"flag"
	"log"
	"os"

	"github.com/eterline/desky-backend/internal/app"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/pkg/logger"
)

var (
	logTimestamp bool
	logPath      string
	genConfig    bool
	configFile   string
)

func init() {

	flag.BoolVar(&logTimestamp, "timestamp", false, `Print starting time in log file name.
With default settings.
True: 'trace.2022.09.07_12.00.00.log'
False: 'trace.log'
	`)

	flag.StringVar(&logPath, "log", "", "Logging file directory.")
	flag.BoolVar(&genConfig, "generate", false, "Generate new config file: 'config.json'")
	flag.StringVar(&configFile, "config", "config.json", "Configuration file path.\nCan be: JSON | YAML | YML extension.")

	flag.Parse()

	if genConfig {
		if err := configuration.GenerateFile("config.json"); err != nil {
			panic(err)
		}
		log.Printf("config file generated: exiting from program")
		os.Exit(0)
	}
}

func main() {
	if err := configuration.Init(configFile); err != nil {
		panic(err)
	}

	conf := configuration.GetConfig()

	logger.InitLogger(
		logger.WithEnv(func() logger.EnvValue {
			if conf.DevMode {
				return logger.DEVELOP
			}
			return logger.PRODUCTION
		}()),
		logger.WithPretty(),
		logger.WithPath(logPath),
	)

	app.Execute()
}
