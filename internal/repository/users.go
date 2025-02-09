package repository

import (
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/repository/storage"
)

type UsersRepository struct {
	DefaultRepository
}

func NewUsersRepository(db *storage.DB) *UsersRepository {
	return &UsersRepository{
		NewDefaultRepository(db),
	}
}

func (r *UsersRepository) All() ([]models.DeskyUserT, error) {

	var users []models.DeskyUserT

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) CreateUser(user *models.DeskyUserT) error {
	return r.db.Create(user).Error
}

func (r *UsersRepository) DeleteUser(id int) error {
	return r.db.Unscoped().Delete(new(models.DeskyUserT), "ID = ?", id).Error
}

func (r *UsersRepository) UserByLogin(login string) (*models.DeskyUserT, error) {

	u := new(models.DeskyUserT)

	if err := r.db.First(u, "Login = ?", login).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UsersRepository) UserById(id int) (*models.DeskyUserT, error) {

	u := new(models.DeskyUserT)

	if err := r.db.First(u, "ID = ?", id).Error; err != nil {
		return nil, err
	}

	return u, nil
}
