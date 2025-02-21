package server

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/repository"
	"github.com/eterline/desky-backend/internal/server/controllers/applications"
	"github.com/eterline/desky-backend/internal/server/controllers/auth"
	"github.com/eterline/desky-backend/internal/server/controllers/exporter"
	"github.com/eterline/desky-backend/internal/server/controllers/frontend"
	"github.com/eterline/desky-backend/internal/server/controllers/monitoring"
	"github.com/eterline/desky-backend/internal/server/controllers/parameters"
	ssh "github.com/eterline/desky-backend/internal/server/controllers/sshlander"
	"github.com/eterline/desky-backend/internal/server/controllers/sys"
	agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"
	"github.com/eterline/desky-backend/internal/services/apps/appsdb"
	"github.com/eterline/desky-backend/internal/services/authorization"
	exporters "github.com/eterline/desky-backend/internal/services/exporter"
	"github.com/eterline/desky-backend/internal/services/router"
	"github.com/eterline/desky-backend/internal/services/system"
	"github.com/eterline/desky-backend/pkg/broker"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/eterline/desky-backend/pkg/storage"
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
	databaseInstance = ctx.Value("sql_database").(*storage.DB)

	f := frontend.Init()

	r.Get("/", router.InitController(f.HTML))
	r.Get("/assets/*", router.InitController(f.Assets))
	r.Get("/static/*", router.InitController(f.Static))
	r.Get("/wallpaper/*", router.InitController(f.WallpaperHandle))

	r.MountWith("/api", api(ctx))

	return
}

// api - setting up api routes
func api(ctx context.Context) (rt *router.RouterService) {
	rt = router.NewRouterService()

	rt.Mount("/apps", controllerApps())
	rt.Mount("/system", controllerSystem(ctx))
	rt.Mount("/agent", controllerAgent(ctx))
	// rt.Mount("/auth", controllerAuth())
	// rt.Mount("/exporter", controllerExporter())
	rt.Mount("/ssh", controllerSSH(ctx))
	rt.Mount("/parameters", controllerParameters())

	return
}

// ================== Setup controller groups ==================

func controllerApps() (routes *router.RouterService) {

	appRepo := repository.NewAppsRepository(databaseInstance)

	srv := applications.Init(appsdb.New(appRepo), appRepo)

	routes = router.MakeSubroute(
		router.NewHandler(router.GET, "/table", srv.ShowTable),
		router.NewHandler(router.POST, "/table/{topic}", srv.CreateApp),

		router.NewHandler(router.DELETE, "/table/{id}", srv.DeleteAppById),
		router.NewHandler(router.DELETE, "/table/{topic}/{number}", srv.DeleteApp),

		router.NewHandler(router.PATCH, "/table/{id}", srv.EditApp),
	)

	return
}

func controllerSystem(ctx context.Context) (routes *router.RouterService) {

	s := sys.Init(ctx, system.New())

	routes = router.MakeSubroute(
		router.NewHandler(router.GET, "/stats", s.Stats),

		router.NewHandler(router.GET, "/systemd", s.SystemdUnits),
		router.NewHandler(router.POST, "/systemd/{unit}/{command}", s.UnitCommand),
	)

	return
}

func controllerAgent(ctx context.Context) (routes *router.RouterService) {

	// agent := ctx.Value("agentmon").(*agentmon.AgentMonitorService)
	broker := ctx.Value("agentmon_mqtt").(*broker.ListenerMQTT)

	agent := agentmon.NewAgentMonitorServiceWithBroker(ctx, broker)
	if err := agent.RunDataUpdater("/agent/stats"); err != nil {
		log.Error(err)
		return
	}

	mon := monitoring.Init(ctx, agent, true)

	routes = router.MakeSubroute(
		router.NewHandler(router.GET, "/monitor", mon.Monitor),
	)

	// routes = router.MakeSubroute(
	// 	router.NewHandler(router.GET, "/monitor", mon.Monitor),
	// )

	return
}

func controllerAuth() (routes *router.RouterService) {

	authService := authorization.New(
		repository.NewUsersRepository(databaseInstance),
	)

	group := auth.Init(authService, authService)

	routes = router.MakeSubroute(
		router.NewHandler(router.POST, "/login", group.Login),
		router.NewHandler(router.POST, "/register", group.Register),

		router.NewHandler(router.GET, "/users", group.Users),
		router.NewHandler(router.DELETE, "/users/{id}", group.Delete),
	)

	return
}

func controllerExporter() (rt *router.RouterService) {

	exporterService := exporters.New(repository.NewExporterRepository(databaseInstance))

	group := exporter.Init(exporterService)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/list", group.ListAll),
		router.NewHandler(router.POST, "/list/{service}", group.Append),
		router.NewHandler(router.DELETE, "/list/{id}", group.Delete),
	)

	return
}

func controllerSSH(ctx context.Context) (routes *router.RouterService) {

	sshRepository := repository.NewSSHLanderRepository(databaseInstance)

	group := ssh.Init(ctx, sshRepository)

	routes = router.MakeSubroute(
		router.NewHandler(router.GET, "/list", group.ListHosts),
		router.NewHandler(router.POST, "/list", group.AppendHost),
		router.NewHandler(router.DELETE, "/list/{id}", group.DeleteHost),

		router.NewHandler(router.GET, "/ping", group.TestHosts),

		router.NewHandler(router.GET, "/connect/{id}", group.ConnectionWS),
	)

	return
}

func controllerParameters() (routes *router.RouterService) {

	coll := logger.NewLoggerCollector()
	logger.HookLevelWriter(coll, logrus.ErrorLevel)

	group := parameters.Init(coll)

	routes = router.MakeSubroute(
		router.NewHandler(router.GET, "/logs", group.GetLogs),
		router.NewHandler(router.GET, "/errors", group.Errors),
	)

	return
}
