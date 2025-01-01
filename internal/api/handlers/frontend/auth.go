package frontend

import "fmt"

type TokenResponse struct {
	Token string `json:"DeskyJWT"`
}

func NewTokenResponse(token string) *TokenResponse {
	return &TokenResponse{
		Token: fmt.Sprintf("Bearer %s", token),
	}
}
