package authorization

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTauthProvider struct {
	SecretKey      []byte
	ExpirationTime time.Duration

	_ struct{}
}

func InitJWT(secret, expirationTime string) *JWTauthProvider {

	timeParsed, err := time.ParseDuration(expirationTime)

	if err != nil {
		timeParsed = time.Hour * 12
	}

	return &JWTauthProvider{
		SecretKey:      []byte(secret),
		ExpirationTime: timeParsed,
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
