package repository

import (
	"github.com/eterline/desky-backend/internal/repository/storage"
)

type DefaultRepository struct {
	db *storage.DB
}

func NewDefaultRepository(db *storage.DB) DefaultRepository {
	return DefaultRepository{
		db: db,
	}
}
