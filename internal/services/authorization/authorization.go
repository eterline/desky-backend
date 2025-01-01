package authorization

import (
	"github.com/eterline/desky-backend/internal/configuration"
)

type AuthForm interface {
	GetPassword() string
	GetUsername() string
}

type Payload interface {
	JSONSting() (string, error)
}

func InitAuth(config *configuration.Configuration) *AuthService {
	return &AuthService{
		JWT: NewJWTauth(config.Server.JWTSecretBytes()),
	}
}

// TODO:
// func (form LoginForm) IsValid() bool {
// 	return false
// }

// auth mocked
func (a *AuthService) IsValid(form AuthForm) bool {
	auth := configuration.GetConfig().Auth

	return form.GetPassword() == auth.Password && form.GetUsername() == auth.Username
}

func (a *AuthService) Token(form Payload) (string, error) {
	out, err := form.JSONSting()
	if err != nil {
		return "", err
	}

	return a.JWT.CreateSignedToken(out)
}
