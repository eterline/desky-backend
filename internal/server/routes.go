package server

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/repository"
	"github.com/eterline/desky-backend/internal/repository/storage"
	"github.com/eterline/desky-backend/internal/server/controllers/applications"
	"github.com/eterline/desky-backend/internal/server/controllers/frontend"
	"github.com/eterline/desky-backend/internal/server/controllers/monitoring"
	"github.com/eterline/desky-backend/internal/server/controllers/sys"
	"github.com/eterline/desky-backend/internal/server/router"
	agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"
	"github.com/eterline/desky-backend/internal/services/apps/appsdb"
	"github.com/eterline/desky-backend/internal/services/system"
	agentclient "github.com/eterline/desky-backend/pkg/agent-client"
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
	rt.Mount("/system", systemControllers())
	rt.Mount("/agent", agentControllers(ctx, c))
	// rt.Mount("/provide", providersControllers())

	return
}

// ================== Setup controller groups ==================

func appsControllers(ctx context.Context) (rt *router.RouterService) {

	db := ctx.Value("database").(*storage.DB)

	repos := repository.NewAppsRepository(db.DB)
	a := appsdb.NewAppService(repos)

	srv := applications.Init(a)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/table", srv.ShowTable),
		router.NewHandler(router.POST, "/table/{topic}", srv.AppendApp),
		router.NewHandler(router.DELETE, "/table/{topic}/{number}", srv.DeleteApp),
	)

	return
}

func systemControllers() (rt *router.RouterService) {

	hinf := system.NewHostInfoService()

	s := sys.Init(hinf)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/info", s.HostInfo),
		router.NewHandler(router.GET, "/systemd", s.SystemdUnits),
		router.NewHandler(router.GET, "/stats", s.HostStatsWS),
	)

	return
}

func agentControllers(ctx context.Context, c *configuration.Configuration) (rt *router.RouterService) {

	agents := agentmon.New(ctx)

	for _, conf := range c.Services.DeskyAgent {

		session, err := agentclient.Reg(conf.API, conf.Token)
		if err != nil {
			log.Errorf("agent init error: %s", err.Error())
		}

		data, ok := session.Info()

		if ok {
			agents.AddSession(session, data.Hostname, data.HostID, conf.API)
			log.Infof("agent init: %s", conf.API)
			continue
		}
		log.Errorf("agent init error: %s", conf.API)
	}

	mon := monitoring.Init(ctx, agents)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/monitor", mon.Monitor),
	)

	return
}

// func providersControllers() (rt *router.RouterService) {

// 	config := configuration.GetConfig()

// 	pvs := providers.Init()

// 	pvs.Register(
// 		providers.PVE,
// 		provider.NewProxmoxProvider(config.Services),
// 	)

// 	rt = router.MakeSubroute(
// 		router.NewHandler(router.GET, "/{service}", pvs.ServiceSessions),
// 		router.NewHandler(router.POST, "/{service}", pvs.ServiceSessions),
// 		router.NewHandler(router.DELETE, "/{service}", pvs.ServiceSessions),
// 	)

// 	return
// }
