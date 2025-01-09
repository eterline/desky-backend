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
			"expire":      pd.expirationDate(),
		},
	)

	return token.SignedString(pd.SecretKey)
}

func (pd *JWTauth) getSecretToken(tokenString *jwt.Token) (interface{}, error) {
	return pd.SecretKey, nil
}

func (pd *JWTauth) expirationDate() int64 {
	return time.Now().Add(pd.ExpirationTime).Unix()
}

func tokenIsExpired(claims jwt.MapClaims) bool {

	exp, ok := claims["expire"]
	if !ok {
		return true
	}

	switch expValue := exp.(type) {

	case float64:
		expTime := time.Unix(int64(expValue), 0)
		return expTime.Before(time.Now())

	case json.Number:
		expInt, err := expValue.Int64()
		if err != nil {
			return true
		}
		expTime := time.Unix(expInt, 0)
		return expTime.Before(time.Now())

	default:
		return true
	}

}

func (pd *JWTauth) TokenIsValid(tokenString string) bool {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, pd.getSecretToken)
	if err != nil {
		return false
	}

	return token.Valid && tokenIsExpired(claims)
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
