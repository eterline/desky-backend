package server

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/repository"
	"github.com/eterline/desky-backend/internal/server/controllers/applications"
	"github.com/eterline/desky-backend/internal/server/controllers/frontend"
	"github.com/eterline/desky-backend/internal/server/controllers/monitoring"
	"github.com/eterline/desky-backend/internal/server/controllers/parameters"
	ssh "github.com/eterline/desky-backend/internal/server/controllers/sshlander"
	"github.com/eterline/desky-backend/internal/server/controllers/sys"
	agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"
	"github.com/eterline/desky-backend/internal/services/apps/appsdb"
	"github.com/eterline/desky-backend/internal/services/router"
	"github.com/eterline/desky-backend/internal/services/system"
	"github.com/eterline/desky-backend/pkg/broker"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/eterline/desky-backend/pkg/storage"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

var (
	log              *logrus.Logger = nil
	databaseInstance *storage.DB    = nil
)

func ConfigRoutes(
	ctx context.Context,
	c *configuration.Configuration,
) (r *router.RouterService) {
	r = router.NewRouterService()

	log = logger.ReturnEntry().Logger
	databaseInstance = ctx.Value(models.DATABASE_CONTEXT_KEY).(*storage.DB)

	f := frontend.Init()

	r.Get("/", router.InitController(f.HTML))
	r.Get("/assets/*", router.InitController(f.Assets))
	r.Get("/static/*", router.InitController(f.Static))
	r.Get("/wallpaper/*", router.InitController(f.WallpaperHandle))

	r.MountWith("/api", api(ctx))

	return
}

// api - setting up api routes
func api(ctx context.Context) *chi.Mux {

	rt := chi.NewRouter()

	rt.Route("/apps", func(r chi.Router) {

		appRepo := repository.NewAppsRepository(databaseInstance)
		srv := applications.Init(appsdb.New(appRepo), appRepo)

		r.Get("/table", router.InitController(srv.ShowTable))
		r.Post("/table/{topic}", router.InitController(srv.CreateApp))
		r.Delete("/table/{id}", router.InitController(srv.DeleteAppById))
		r.Patch("/table/{id}", router.InitController(srv.EditApp))
	})

	rt.Route("/system", func(r chi.Router) {

		srv := sys.Init(ctx, system.New())

		r.Get("/stats", router.InitController(srv.Stats))
		r.Get("/systemd", router.InitController(srv.SystemdUnits))
	})

	rt.Route("/agent", func(r chi.Router) {

		broker := ctx.Value(models.MESSAGE_BROKER_CONTEXT_KEY).(*broker.ListenerMQTT)
		agent := agentmon.NewAgentMonitorServiceWithBroker(ctx, broker)
		if err := agent.RunDataUpdater("/agent/stats"); err != nil {
			log.Error(err)
			return
		}
		mon := monitoring.Init(ctx, agent, true)

		r.Get("/monitor", router.InitController(mon.Monitor))
	})

	rt.Route("/ssh", func(r chi.Router) {

		sshRepository := repository.NewSSHLanderRepository(databaseInstance)
		srv := ssh.Init(ctx, sshRepository)

		r.Get("/list", router.InitController(srv.ListHosts))
		r.Post("/list", router.InitController(srv.AppendHost))
		r.Delete("/list/{id}", router.InitController(srv.DeleteHost))

		r.Get("/ping", router.InitController(srv.TestHosts))
		r.Get("/connect/{id}", router.InitController(srv.ConnectionWS))
	})

	rt.Route("/parameters", func(r chi.Router) {
		coll := logger.NewLoggerCollector()
		logger.HookLevelWriter(coll, logrus.ErrorLevel)
		srv := parameters.Init(coll)

		r.Get("/logs", router.InitController(srv.GetLogs))
		r.Get("/errors", router.InitController(srv.Errors))
	})

	return rt
}
