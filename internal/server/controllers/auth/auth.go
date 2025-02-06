package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/router/handler"
	"github.com/eterline/desky-backend/internal/utils"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

type AuthProvder interface {
	All() ([]*models.DeskyUser, error)
	Register(*models.DeskyUser) error
	Edit(user *models.DeskyUser, id int) error
	Delete(id int) error
}

type VerifyProvider interface {
	Verify(login, password string) (*models.DeskyUser, error)
	VerifyWithID(id int, login, password string) (*models.DeskyUser, error)
}

var log *logrus.Logger = nil

type AuthHandlerGroup struct {
	service    AuthProvder
	verificate VerifyProvider
}

func Init(auth AuthProvder, verif VerifyProvider) *AuthHandlerGroup {
	log = logger.ReturnEntry().Logger
	return &AuthHandlerGroup{
		service:    auth,
		verificate: verif,
	}
}

func (au *AuthHandlerGroup) Register(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "auth.register"

	user := new(models.DeskyUser)

	if err := handler.DecodeRequest(r, user); err != nil {
		fmt.Println(err)
		return op, err
	}
	log.Debugf("register request user: %s", utils.PrettyString(user))

	if err := au.service.Register(user); err != nil {
		return op, handler.NewErrorResponse(
			http.StatusNotImplemented,
			errors.New("register failed"),
		)
	}

	return op, handler.StatusOK(w, "ok")
}

func (au *AuthHandlerGroup) Login(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "auth.login"

	user := new(models.DeskyUser)

	if err := handler.DecodeRequest(r, user); err != nil {
		return op, handler.StatusUnauthorized()
	}
	log.Debugf("login request user: %s", utils.PrettyString(user))

	_, err = au.verificate.Verify(user.Login, user.Password)
	if err != nil {
		log.Error(err)
		return op, handler.UnauthorizedErrorResponse()
	}

	return op, nil
}

func (au *AuthHandlerGroup) Delete(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "auth.delete"
	user := new(models.DeskyUser)

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	if err := handler.DecodeRequest(r, user); err != nil {
		return op, err
	}
	log.Debugf("delete user request: %s", utils.PrettyString(user))

	if _, err := au.verificate.VerifyWithID(q.GetInt("id"), user.Login, user.Password); err != nil {
		log.Error(err)
		return op, handler.UnauthorizedErrorResponse()
	}

	if err := au.service.Delete(q.GetInt("id")); err != nil {
		log.Error(err)
		return op, handler.NewErrorResponse(
			http.StatusNotImplemented,
			errors.New("user not found"),
		)
	}

	return op, handler.StatusOK(w, "user deleted")
}

func (au *AuthHandlerGroup) Users(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "auth.users"

	data, err := au.service.All()
	if err != nil {
		log.Error(err)
		return op, err
	}

	if len(data) < 1 {
		return op, handler.NoContentResponse()
	}

	return op, handler.WriteJSON(w, http.StatusOK, data)
}
