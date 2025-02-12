package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/pkg/storage"
)

var (
	dbFile      string
	migrateList []any
)

func init() {
	flag.BoolFunc("gen", "To generate configuration file.", genConfig)
	flag.StringVar(&dbFile, "file", "desky.db", "Set up database file for migration")
	flag.Parse()

	migrateList = []any{
		new(models.AppsTopicT),
		new(models.AppsInstancesT),
		new(models.DeskyUserT),
		new(models.ExporterInfoT),
		new(models.SSHCredentialsT),
		new(models.SSHSystemTypesT),
		new(models.SSHSecureT),
	}
}

func main() {

	db := storage.New(storage.NewStorageSQLite("desky.db"), nil)

	err := db.Connect()
	defer db.Close()
	if err != nil {
		panic(fmt.Sprintf("db connect error: %v", err))
	}

	if err := db.MigrateTables(migrateList...); err != nil {
		panic(fmt.Sprintf("db migration error: %v", err))
	}

	fmt.Printf("db migrated to: %s \n", db.Source())
	fmt.Println("exiting migrator")
}

func genConfig(string) error {
	if err := configuration.Migrate(configuration.FileName, 0644); err != nil {
		panic(err)
	}
	fmt.Println("Migration: default config generated")
	os.Exit(0)
	return nil
}
