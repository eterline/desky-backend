package repository

import "gorm.io/gorm"

type DefaultRepository struct {
	db *gorm.DB
}
