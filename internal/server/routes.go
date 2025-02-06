package server

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/repository"
	"github.com/eterline/desky-backend/internal/repository/storage"
	"github.com/eterline/desky-backend/internal/server/controllers/applications"
	"github.com/eterline/desky-backend/internal/server/controllers/auth"
	"github.com/eterline/desky-backend/internal/server/controllers/exporter"
	"github.com/eterline/desky-backend/internal/server/controllers/frontend"
	"github.com/eterline/desky-backend/internal/server/controllers/monitoring"
	"github.com/eterline/desky-backend/internal/server/controllers/sys"
	agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"
	"github.com/eterline/desky-backend/internal/services/apps/appsdb"
	"github.com/eterline/desky-backend/internal/services/authorization"
	exporters "github.com/eterline/desky-backend/internal/services/exporter"
	"github.com/eterline/desky-backend/internal/services/router"
	"github.com/eterline/desky-backend/internal/services/system"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func ConfigRoutes(
	ctx context.Context,
	c *configuration.Configuration,
) (r *router.RouterService) {
	r = router.NewRouterService()

	log = logger.ReturnEntry().Logger

	f := frontend.Init()

	r.Get("/", router.InitController(f.HTML))
	r.Get("/assets/*", router.InitController(f.Assets))
	r.Get("/static/*", router.InitController(f.Static))
	r.Get("/wallpaper/*", router.InitController(f.WallpaperHandle))

	r.MountWith("/api", api(ctx, c))

	return
}

// api - setting up api routes
func api(
	ctx context.Context,
	c *configuration.Configuration,
) (rt *router.RouterService) {
	rt = router.NewRouterService()

	rt.Mount("/apps", appsControllers(ctx))
	rt.Mount("/system", systemControllers(ctx))
	rt.Mount("/agent", agentControllers(ctx, c))
	rt.Mount("/auth", authControllers(ctx, c))
	rt.Mount("/exporter", exporterControllers(ctx))

	return
}

// ================== Setup controller groups ==================

func appsControllers(ctx context.Context) (rt *router.RouterService) {

	db := ctx.Value("db").(*storage.DB)

	srv := applications.Init(
		appsdb.NewAppService(repository.NewAppsRepository(db.DB)),
	)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/table", srv.ShowTable),
		router.NewHandler(router.POST, "/table/{topic}", srv.AppendApp),
		router.NewHandler(router.DELETE, "/table/{topic}/{number}", srv.DeleteApp),
	)

	return
}

func systemControllers(ctx context.Context) (rt *router.RouterService) {

	hinf := system.NewHostInfoService()

	s := sys.Init(ctx, hinf)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/stats", s.Stats),

		router.NewHandler(router.GET, "/systemd", s.SystemdUnits),
		router.NewHandler(router.POST, "/systemd/{unit}/{command}", s.UnitCommand),
	)

	return
}

func agentControllers(ctx context.Context, c *configuration.Configuration) (rt *router.RouterService) {

	a := ctx.Value("agentmon").(*agentmon.AgentMonitorService)
	mon := monitoring.Init(ctx, a, true)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/monitor", mon.Monitor),
	)

	return
}

func authControllers(ctx context.Context, c *configuration.Configuration) (rt *router.RouterService) {

	db := ctx.Value("db").(*storage.DB)

	authService := authorization.New(
		repository.NewUsersRepository(db.DB),
	)

	authConrollers := auth.Init(authService, authService)

	rt = router.MakeSubroute(
		router.NewHandler(router.POST, "/login", authConrollers.Login),
		router.NewHandler(router.POST, "/register", authConrollers.Register),

		router.NewHandler(router.GET, "/users", authConrollers.Users),
		router.NewHandler(router.DELETE, "/users/{id}", authConrollers.Delete),
	)

	return
}

func exporterControllers(ctx context.Context) (rt *router.RouterService) {

	db := ctx.Value("db").(*storage.DB)

	service := exporters.NewExporterService(
		repository.NewExporterRepository(db.DB),
	)

	exc := exporter.Init(service)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/exporter", exc.ListAll),
		router.NewHandler(router.POST, "/exporter/{service}", exc.Append),
		router.NewHandler(router.DELETE, "/exporter/{id}", exc.Delete),
	)

	return
}
