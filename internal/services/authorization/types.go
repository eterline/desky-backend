package authorization

import (
	"net/http"
	"time"

	"github.com/eterline/desky-backend/internal/api/handlers"
)

// Implements JWT authentication for web-server
type JWTauth struct {
	SecretKey      []byte
	ExpirationTime time.Duration

	_ struct{}
}

type AuthService struct {
	JWT *JWTauth
}

type LoginForm struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (f *LoginForm) GetPassword() string {
	return f.Password
}

func (f *LoginForm) GetUsername() string {
	return f.Username
}

type AuthorizedResponse struct {
	Token string `json:"DeskyJWT"`
}

func DecodeCredentials(r *http.Request) (*LoginForm, error) {
	form := &LoginForm{}

	if err := handlers.DecodeRequest(r, form); err != nil {
		return nil, err
	}

	return form, nil
}
