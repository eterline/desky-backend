package application

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/server"
	agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"
	agentclient "github.com/eterline/desky-backend/pkg/agent-client"
	"github.com/sirupsen/logrus"
)

func Exec(
	ctx context.Context,
	log *logrus.Logger,
	config *configuration.Configuration,
) {
	// ================= Services additional =================

	mon := agentmon.New(ctx)

	for _, agent := range config.Services.DeskyAgent {
		log.Infof("connecting to agent: %s", agent.API)
		cl, err := agentclient.Reg(agent.API, agent.Token)

		if err != nil {
			log.Errorf("skip: %s. error: %s", agent.API, err.Error())
			continue
		}

		log.Infof(
			"agent successfully connected. hostname: %s. id: %s",
			cl.Info.Hostname, cl.Info.HostID,
		)

		mon.AddSession(cl, cl.Info.Hostname, cl.Info.HostID, agent.API)
	}

	ctx = context.WithValue(ctx, "agentmon", mon)

	// ================= Server parameters =================

	srv := server.New(
		config.ServerSocket(),
		config.SSL().CertFile,
		config.SSL().KeyFile,
		config.Server.Name,
	)

	router := server.ConfigRoutes(ctx, config)

	router.Get("/config", PreferencesHandler(false, true))
	router.Get("/health", HealthHandler)

	srv.Router(router)

	// ================= Run main server parameters =================

	go func() {
		log.Infof("server start at: %s", config.URLString())
		defer log.Info("server closed")

		if err := srv.Run(config.SSL().TLS); err != nil {
			log.Errorf("server running error: %s", err)
		}
	}()

	<-ctx.Done()

	if err := srv.Stop(); err != nil {
		log.Errorf("server close error: %s", err.Error())
	}
}
