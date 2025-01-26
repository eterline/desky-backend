package server

import (
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/server/controllers/applications"
	"github.com/eterline/desky-backend/internal/server/controllers/frontend"
	"github.com/eterline/desky-backend/internal/server/controllers/providers"
	"github.com/eterline/desky-backend/internal/server/controllers/sys"
	"github.com/eterline/desky-backend/internal/server/router"
	"github.com/eterline/desky-backend/internal/services/apps/appsfile"
	"github.com/eterline/desky-backend/internal/services/provider"
	"github.com/eterline/desky-backend/internal/services/system"
)

// public - setting up public routes
func public(rt *router.RouterService) {
	f := frontend.Init()

	rt.Get("/", router.InitController(f.HTML))
	rt.Get("/config", router.InitController(f.Assets))

	rt.Get("/assets/*", router.InitController(f.Assets))
	rt.Get("/static/*", router.InitController(f.Static))
	rt.Get("/wallpaper/*", router.InitController(f.WallpaperHandle))
}

// api - setting up api routes
func api() (rt *router.RouterService) {
	rt = router.NewRouterService()

	rt.Mount("/apps", appsControllers())
	rt.Mount("/system", systemControllers())
	rt.Mount("/provide", providersControllers())

	return
}

// ================== Setup controller groups ==================

func appsControllers() (rt *router.RouterService) {

	a, err := appsfile.Init("./apps.json")
	if err != nil {
		panic(err)
	}

	srv := applications.Init(a)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/table", srv.ShowTable),
		router.NewHandler(router.POST, "/table/{topic}", srv.AppendApp),
		router.NewHandler(router.DELETE, "/table/{topic}/{number}", srv.DeleteApp),
	)

	log.Debug("apps controllers registered")

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

	log.Debug("system controllers registered")

	return
}

func providersControllers() (rt *router.RouterService) {

	config := configuration.GetConfig()

	pvs := providers.Init()

	pvs.Register(
		providers.PVE,
		provider.NewProxmoxProvider(config.Services),
	)

	rt = router.MakeSubroute(
		router.NewHandler(router.GET, "/{service}", pvs.ServiceSessions),
		router.NewHandler(router.POST, "/{service}", pvs.ServiceSessions),
		router.NewHandler(router.DELETE, "/{service}", pvs.ServiceSessions),
	)

	log.Debug("providers controllers registered")

	return
}
