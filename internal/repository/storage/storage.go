package storage

import (
	"context"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DefaultName = "desky.db"

type DB struct {
	name   string
	config *gorm.Config

	*gorm.DB
}

func New(file string, logger logger.Interface) *DB {
	return &DB{
		name: file,
		DB:   nil,
		config: &gorm.Config{
			Logger: logger,
		},
	}
}

func (db *DB) Test() bool {

	_, err := os.Stat(db.name)

	switch {
	case os.IsNotExist(err) == true:
		db.name = DefaultName
		return false

	case err == nil:
		return true

	default:
		panic(err)
	}
}

func (db *DB) Source() string {
	return db.name
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

func (db *DB) MigrateTables(tables ...any) error {
	base, err := gorm.Open(sqlite.Open(db.name), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := base.AutoMigrate(tables...); err != nil {
		return err
	}

	return nil
}
