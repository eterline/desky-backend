package authorization

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWTExpirationTime time.Duration = time.Hour * 12

func ConvertSecret(v string) []byte {
	return []byte(v)
}

func NewJWTauth(secret []byte) *JWTauth {
	return &JWTauth{
		SecretKey:      secret,
		ExpirationTime: JWTExpirationTime,
	}
}

func (pd *JWTauth) CreateSignedToken(credentials string) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"credentials": credentials,
			"expiration":  pd.ExpirationTime,
		},
	)

	return token.SignedString(pd.SecretKey)
}

func (pd *JWTauth) TokenIsValid(tokenString string) bool {

	var getSecretToken = func(tokenString *jwt.Token) (interface{}, error) {
		return pd.SecretKey, nil
	}

	token, err := jwt.Parse(tokenString, getSecretToken)
	if err != nil {
		return false
	}

	return token.Valid
}

type JWTPayload struct {
	Username string    `json:"username"`
	Time     time.Time `json:"time"`
}

func NewPayload(username string) *JWTPayload {
	return &JWTPayload{
		Username: username,
		Time:     time.Now(),
	}
}

func (f *JWTPayload) JSONSting() (string, error) {
	s, err := json.Marshal(f)
	return string(s), err
}
