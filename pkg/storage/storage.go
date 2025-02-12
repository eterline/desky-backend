package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DefaultName = "desky.db"

type DB struct {
	provide StorageProvider
	config  *gorm.Config
	*gorm.DB
}

func New(provider StorageProvider, logger logger.Interface) *DB {
	return &DB{
		provide: provider,
		config: &gorm.Config{
			Logger: logger,
		},
		DB: nil,
	}
}

func (db *DB) Source() string {
	return fmt.Sprintf("%s=%s", db.provide.StorageType(), db.provide.Source())
}

func (db *DB) Connect() error {
	base, err := gorm.Open(db.provide.Socket(), db.config)
	if err != nil {
		return err
	}

	base.Logger.Info(
		context.Background(), "db connected: %s=%s",
		db.provide.StorageType(), db.provide.Source(),
	)

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
	base, err := gorm.Open(db.provide.Socket(), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := base.AutoMigrate(tables...); err != nil {
		return err
	}

	return nil
}
