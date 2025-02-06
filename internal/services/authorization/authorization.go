package authorization

import (
	"errors"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/hash"
	"gorm.io/gorm"
)

type UserRepository interface {
	All() ([]models.DeskyUserT, error)
	CreateUser(user *models.DeskyUserT) error
	DeleteUser(id int) error
	UserByLogin(login string) (*models.DeskyUserT, error)
	UserById(id int) (*models.DeskyUserT, error)
}

type AuthorizationService struct {
	repository UserRepository
	hash       *hash.HashService
}

func New(r UserRepository) *AuthorizationService {
	return &AuthorizationService{
		repository: r,
		hash:       hash.NewHashService(hash.SHA512, []byte("random")),
	}
}

func (aus *AuthorizationService) All() ([]*models.DeskyUser, error) {
	users, err := aus.repository.All()
	if err != nil {
		return nil, err
	}

	data := make([]*models.DeskyUser, len(users))

	for i, usr := range users {
		data[i] = models.NewDeskyUser(usr.ID, usr.Login, "")
	}

	return data, nil
}

func (aus *AuthorizationService) Register(user *models.DeskyUser) error {

	tReq, err := aus.repository.UserByLogin(user.Login)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err != gorm.ErrRecordNotFound && tReq.Login == user.Login {
		return errors.New("user already exists")
	}

	hashedPassword := aus.hash.StringHash(user.Password)
	userT := models.NewDeskyUserT(user.Login, hashedPassword.String())

	return aus.repository.CreateUser(userT)
}

func (aus *AuthorizationService) Edit(user *models.DeskyUser, id int) error {
	panic("")
}

func (aus *AuthorizationService) Delete(id int) error {
	return aus.repository.DeleteUser(id)
}

func (aus *AuthorizationService) VerifyWithID(id int, login, password string) (*models.DeskyUser, error) {

	user, err := aus.repository.UserById(id)
	if err != nil {
		return nil, err
	}

	if !(aus.hash.EqStrings(user.Password, password) && int(user.ID) == id) {
		return nil, ErrVerifyPassword
	}

	return models.NewDeskyUser(user.ID, user.Login, user.Password), nil
}

func (aus *AuthorizationService) Verify(login, password string) (*models.DeskyUser, error) {

	user, err := aus.repository.UserByLogin(login)
	if err != nil {
		return nil, err
	}

	ok := aus.hash.EqStrings(user.Password, password)
	if !ok {
		return nil, ErrVerifyPassword
	}

	return models.NewDeskyUser(user.ID, user.Login, user.Password), nil
}
