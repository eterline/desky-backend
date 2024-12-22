package app

import (
	"github.com/eterline/desky-backend/internal/api"
	"github.com/sirupsen/logrus"
)

type App struct {
	Server *api.Server
	Log    *logrus.Logger
	// Services
}
