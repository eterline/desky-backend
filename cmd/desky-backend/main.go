package main

import (
	"flag"
	"os"

	"github.com/eterline/desky-backend/internal/app"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/pkg/logger"
)

var (
	logConfig        logger.LoggingConfig
	isGenerateConfig bool
	configFile       string
)

func init() {

	flag.BoolVar(&logConfig.Renew, "timestamp", false, `Print starting time in log file name.
With default settings.
True: 'trace.2022.09.07_12.00.00.log'
False: 'trace.log'
	`)

	flag.StringVar(&logConfig.Dir, "log", "", "Logging file directory.")
	flag.UintVar(&logConfig.Level, "loglvl", 0, "Logging level.")

	flag.BoolVar(&isGenerateConfig, "generate", false, "Generate new config file: 'config.json'")
	flag.StringVar(&configFile, "config", "config.json", "Configuration file path.\nCan be: JSON | YAML | YML extension.")

	flag.Parse()

	logger.InitWithConfig(logConfig)
}

func main() {
	log := logger.ReturnEntry()

	if isGenerateConfig {
		if err := configuration.GenerateFile("config.json"); err != nil {
			log.Fatalf("can't generate config file: %s", err.Error())
		}
		log.Info("config file generated: exiting from program")
		os.Exit(0)
	}

	if err := configuration.Init(configFile); err != nil {
		log.Errorf("failed init config file: %s", err.Error())
	}

	app.Execute()
}
