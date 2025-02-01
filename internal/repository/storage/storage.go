package storage

import (
	"context"

	"github.com/eterline/desky-backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DefaultName = "desky.db"

type DB struct {
	name   string
	DB     *gorm.DB
	config *gorm.Config
}

func TestFile(file string) bool {
	return file != ""
}

func New(name string) *DB {
	return &DB{
		name: name,
		DB:   nil,
		config: &gorm.Config{
			Logger: NewLog(),
		},
	}
}

func (db *DB) Connect() error {
	base, err := gorm.Open(sqlite.Open(db.name), db.config)
	if err != nil {
		return err
	}
	db.DB = base
	return nil
}

func (db *DB) Close() error {

	db.DB.Logger.Info(context.Background(), "closing db connection")

	dbInstance, err := db.DB.DB()
	if err != nil {
		panic(err)
	}

	return dbInstance.Close()
}

func (db *DB) MigrateTables() error {
	base, err := gorm.Open(sqlite.Open(db.name), &gorm.Config{})
	if err != nil {
		return err
	}

	err = base.AutoMigrate(
		&models.AppsTopicT{},
		&models.AppsInstancesT{},
		&models.WidgetsT{},
	)

	if err != nil {
		return err
	}

	return nil
}
