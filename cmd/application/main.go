package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/pkg/logger"
)

func init() {
	flag.BoolFunc("gen", "to generate configuration file", func(s string) error {
		if err := configuration.GenerateFile("config.yaml"); err != nil {
			panic(err)
		}
		os.Exit(0)
		return nil
	})

	flag.Parse()

	configuration.MustConfig("config.yaml")
}

// @title		Desky API test
// @version	1.0
// @BasePath	/api/v1
func main() {
	if err := logger.InitLogger(
		logger.WithPretty(),
		logger.WithPath(""),
	); err != nil {
		log.Println(err)
	}

	c := configuration.GetConfig()

	fmt.Println(c.SSL().CertFile)
}
