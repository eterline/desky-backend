package repository

import (
	"github.com/eterline/desky-backend/pkg/storage"
)

type DefaultRepository struct {
	db *storage.DB
}

func NewDefaultRepository(db *storage.DB) DefaultRepository {
	return DefaultRepository{
		db: db,
	}
}
