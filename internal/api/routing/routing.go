package routing

import (
	"fmt"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/api/handlers/applications"
	"github.com/eterline/desky-backend/internal/api/handlers/frontend"
	"github.com/eterline/desky-backend/internal/api/handlers/proxmox"
	"github.com/eterline/desky-backend/internal/api/handlers/sys"
	mws "github.com/eterline/desky-backend/internal/api/middlewares"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/services/authorization"
	"github.com/eterline/desky-backend/internal/services/cache"
	"github.com/eterline/desky-backend/internal/services/system"
	"github.com/go-chi/chi"
)

type APIRouting struct {
	BasePath string
	MW       *mws.MiddleWare
}

func InitAPIRouting(ver int) *APIRouting {
	return &APIRouting{
		BasePath: fmt.Sprintf("/api/v%d", ver),
		MW:       mws.Init(),
	}
}

func (rt *APIRouting) pathWithBase(pt string) string {
	return rt.BasePath + pt
}

func (rt *APIRouting) ConfigRoutes() *chi.Mux {

	router := rt.setBaseRouting()
	protectedRoute := chi.NewRouter()

	protectedRoute.Use(rt.MW.AuthorizationJWT)

	// controller for auth status checking
	protectedRoute.Get("/check", handlers.InitController(frontend.AccessCheck))

	protectedRoute.Mount("/apps", setApplicationsRouting())
	protectedRoute.Mount("/pve", setProxmoxRouting())
	protectedRoute.Mount("/system", setSystemRouting())

	router.Mount(rt.BasePath, protectedRoute)

	return router
}

func (rt *APIRouting) setBaseRouting() *chi.Mux {
	chi := chi.NewMux()

	chi.Use(rt.MW.PanicRecoverer, rt.MW.Logging, rt.MW.Compressor)

	front := frontend.Init(authorization.InitAuth(configuration.GetConfig()))

	chi.Post("/login", handlers.InitController(front.Login))

	chi.Get("/", handlers.InitController(front.HTML))
	chi.Get("/assets/*", handlers.InitController(front.Assets))
	chi.Get("/static/*", handlers.InitController(front.Static))
	chi.Get("/wallpaper/*", handlers.InitController(front.WallpaperHandle))

	return chi
}

func setApplicationsRouting() *chi.Mux {
	as := applications.Init("apps.json")

	return BuildSubroute(
		RoutesConfig{
			HandlerParam{"GET", "/table", as.ReturnAppsTable},
			HandlerParam{"POST", "/table/{topic}", as.AppendApp},
			HandlerParam{"DELETE", "/table/{topic}/{number}", as.DeleteApp},
		},
	)
}

func setProxmoxRouting() *chi.Mux {
	pve := proxmox.Init(cache.GetEntry())

	return BuildSubroute(
		RoutesConfig{
			HandlerParam{"GET", "/sessions", pve.Sessions},
			HandlerParam{"GET", "/{session}/{node}/status", pve.NodeStatus},
			HandlerParam{"GET", "/{session}/{node}/devices", pve.DeviceList},

			HandlerParam{"GET", "/{session}/{node}/apt/updates", pve.AptUpdates},
			HandlerParam{"POST", "/{session}/{node}/apt/update", pve.AptUpdate},

			HandlerParam{"POST", "/{session}/{node}/devices/{vmid}/{command}", pve.DeviceCommand},
		},
	)
}

func setSystemRouting() *chi.Mux {
	sys := sys.Init(system.NewHostInfoService(), cache.GetEntry())

	return BuildSubroute(
		RoutesConfig{
			HandlerParam{"GET", "/info", sys.HostInfo},

			HandlerParam{"GET", "/systemd/status", sys.SystemdUnits},
			HandlerParam{"POST", "/systemd/{unit}/{command}", sys.UnitCommand},

			HandlerParam{"GET", "/stats", sys.HostStatsWS},
			HandlerParam{"GET", "/tty", sys.TtyWS},
		},
	)
}
