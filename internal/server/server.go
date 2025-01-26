package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/server/router"
	"github.com/eterline/desky-backend/internal/utils"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = nil

func Setup(config *configuration.Configuration) *http.Server {
	log = logger.ReturnEntry().Logger
	defer logStart(config)

	tls := func() *tls.Config {
		if config.SSL().TLS {
			return &tls.Config{
				ServerName: config.Server.Name,
			}
		}
		return nil
	}()

	return &http.Server{
		Addr:      config.ServerSocket(),
		Handler:   ConfigRoutes(),
		TLSConfig: tls,
	}
}

func ConfigRoutes() (r *router.RouterService) {
	r = router.NewRouterService()

	public(r)
	log.Info("public routes configured")

	r.MountWith("/api", api())
	log.Info("api routes configured")

	return
}

func logStart(config *configuration.Configuration) {
	log.Debugf("init server with config: %s", utils.PrettyString(config.Server))
	log.Infof("https mode: %v", config.SSL().TLS)

	var addr string = config.ServerSocket()

	if config.SSL().TLS {
		addr = fmt.Sprintf("https://%s", addr)
	} else {
		addr = fmt.Sprintf("http://%s", addr)
	}

	log.Infof("server listen on: %s", addr)
}
