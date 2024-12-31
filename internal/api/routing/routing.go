package routing

import (
	"fmt"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/api/handlers/applications"
	"github.com/eterline/desky-backend/internal/api/handlers/frontend"
	"github.com/eterline/desky-backend/internal/api/handlers/proxmox"
	"github.com/eterline/desky-backend/internal/api/handlers/sys"
	mws "github.com/eterline/desky-backend/internal/api/middlewares"
	"github.com/go-chi/chi"
)

type APIRouting struct {
	BasePath string
}

func InitAPIRouting(ver int) *APIRouting {
	return &APIRouting{
		BasePath: fmt.Sprintf("/api/v%d", ver),
	}
}

func (rt *APIRouting) pathWithBase(pt string) string {
	return rt.BasePath + pt
}

func (rt *APIRouting) ConfigRoutes() *chi.Mux {
	router := setBaseRouting()

	router.Mount(rt.pathWithBase("/apps"), setApplicationsRouting())
	router.Mount(rt.pathWithBase("/pve"), setProxmoxRouting())
	router.Mount(rt.pathWithBase("/system"), setSystemRouting())

	return router
}

func setBaseRouting() *chi.Mux {
	chi := chi.NewMux()
	mw := mws.Init()

	chi.Use(mw.Logging)

	front := frontend.Init()

	chi.Get("/", handlers.InitController(front.HTML))
	chi.Get("/assets/*", handlers.InitController(front.Assets))
	chi.Get("/static/*", handlers.InitController(front.Static))

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
	pve := proxmox.Init()

	return BuildSubroute(
		RoutesConfig{
			HandlerParam{"GET", "/{session}/{node}/status", pve.NodeStatus},

			HandlerParam{"GET", "/{session}/{node}/devices", pve.DeviceList},

			HandlerParam{"POST", "/{session}/{node}/devices/{vmid}/start", pve.DeviceStart},
			HandlerParam{"POST", "/{session}/{node}/devices/{vmid}/shutdown", pve.DeviceShutdown},
			HandlerParam{"POST", "/{session}/{node}/devices/{vmid}/stop", pve.DeviceStop},
			HandlerParam{"POST", "/{session}/{node}/devices/{vmid}/suspend", pve.DeviceSuspend},
			HandlerParam{"POST", "/{session}/{node}/devices/{vmid}/resume", pve.DeviceResume},
		},
	)
}

func setSystemRouting() *chi.Mux {
	sys := sys.Init()

	return BuildSubroute(
		RoutesConfig{
			HandlerParam{"GET", "/info", sys.HostInfo},
			HandlerParam{"GET", "/stats", sys.HostStatsWS},
			HandlerParam{"GET", "/tty", sys.TtyWS},

			HandlerParam{"GET", "/systemd/status", sys.SystemdUnits},
		},
	)
}
