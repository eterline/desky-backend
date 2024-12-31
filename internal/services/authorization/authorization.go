package authorization

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWTExpirationTime time.Duration = time.Hour * 12

// Implements JWT authentication for web-server
type JWTauthProvider struct {
	SecretKey      []byte
	ExpirationTime time.Duration

	_ struct{}
}

func ConvertSecret(v string) []byte {
	return []byte(v)
}

func NewJWTauthProvider(secret []byte) *JWTauthProvider {
	return &JWTauthProvider{
		SecretKey:      secret,
		ExpirationTime: JWTExpirationTime,
	}
}

func (pd *JWTauthProvider) CreateSignedToken(credentials string) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"credentials": credentials,
			"expiration":  pd.ExpirationTime,
		},
	)

	return token.SignedString(pd.SecretKey)
}

func (pd *JWTauthProvider) TokenIsValid(tokenString string) bool {

	var getSecretToken = func(tokenString *jwt.Token) (interface{}, error) {
		return pd.SecretKey, nil
	}

	token, err := jwt.Parse(tokenString, getSecretToken)
	if err != nil {
		return false
	}

	return token.Valid
}
