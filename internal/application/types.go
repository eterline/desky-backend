package application

import "net/http"

type Application struct {
	Server *http.Server
}
