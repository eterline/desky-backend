package application

import "net/http"

type Application struct {
	Server *http.Server
}

type PresencesResponse struct {
	DarkTheme  bool   `json:"dark-theme"`
	Language   string `json:"language"`
	Background string `json:"background"`
	Auth       bool   `json:"auth"`
}
