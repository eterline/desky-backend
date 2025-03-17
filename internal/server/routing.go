package server

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/repository"
	"github.com/eterline/desky-backend/internal/server/controllers"
	middlewares "github.com/eterline/desky-backend/internal/server/middleware"
	agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"
	"github.com/eterline/desky-backend/internal/services/apps/appsdb"
	"github.com/eterline/desky-backend/internal/services/handler"
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

func ConfigRoutes(ctx context.Context, c *configuration.Configuration) (r *chi.Mux) {

	r = chi.NewRouter()

	log = logger.ReturnEntry().Logger
	databaseInstance = ctx.Value(models.DATABASE_CONTEXT_KEY).(*storage.DB)

	f := controllers.InitFronEnd()

	r.Get("/", handler.InitController(f.HTML))
	r.Get("/assets/*", handler.InitController(f.Assets))
	r.Get("/static/*", handler.InitController(f.Static))
	r.Get("/wallpaper/*", handler.InitController(f.WallpaperHandle))

	r.With(
		middlewares.CorsPolicy,
		middlewares.FilterContentType,
		middlewares.PreSetHeaders,
	).Mount("/api", api(ctx))

	return
}

// api - setting up api routes
func api(ctx context.Context) *chi.Mux {

	rt := chi.NewRouter()

	rt.Route("/apps", func(r chi.Router) {

		appRepo := repository.NewAppsRepository(databaseInstance)
		srv := controllers.InitApplications(appsdb.New(appRepo), appRepo)

		r.Get("/table", handler.InitController(srv.ShowTable))
		r.Post("/table/{topic}", handler.InitController(srv.CreateApp))
		r.Delete("/table/{id}", handler.InitController(srv.DeleteAppById))
		r.Patch("/table/{id}", handler.InitController(srv.EditApp))
	})

	rt.Route("/system", func(r chi.Router) {

		srv := controllers.InitSystem(ctx, system.New())

		r.Get("/stats", handler.InitController(srv.Stats))
		r.Get("/systemd", handler.InitController(srv.SystemdUnits))
	})

	rt.Route("/agent", func(r chi.Router) {

		broker := ctx.Value(models.MESSAGE_BROKER_CONTEXT_KEY).(*broker.ListenerMQTT)
		agent := agentmon.NewAgentMonitorServiceWithBroker(ctx, broker)
		if err := agent.RunDataUpdater("/agent/stats"); err != nil {
			log.Error(err)
			return
		}
		mon := controllers.InitMonitoring(ctx, agent, true)

		r.Get("/monitor", handler.InitController(mon.Monitor))
	})

	rt.Route("/ssh", func(r chi.Router) {

		sshRepository := repository.NewSSHLanderRepository(databaseInstance)
		srv := controllers.InitSSHlander(ctx, sshRepository)

		r.Get("/list", handler.InitController(srv.ListHosts))
		r.Post("/list", handler.InitController(srv.AppendHost))
		r.Delete("/list/{id}", handler.InitController(srv.DeleteHost))

		r.Get("/ping", handler.InitController(srv.TestHosts))
		r.Get("/connect/{id}", handler.InitController(srv.ConnectionWS))
	})

	rt.Route("/parameters", func(r chi.Router) {
		coll := logger.NewLoggerCollector()
		logger.HookLevelWriter(coll, logrus.ErrorLevel)
		srv := controllers.InitParameters(coll)

		r.Get("/logs", handler.InitController(srv.GetLogs))
		r.Get("/errors", handler.InitController(srv.Errors))
	})

	return rt
}
