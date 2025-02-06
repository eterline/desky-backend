package router

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/services/router/handler"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = nil

const (
	OPTIONS HandlerMethod = iota
	GET
	HEAD
	POST
	PUT
	PATCH
	DELETE
	TRACE
)

func NewRouterService() *RouterService {

	log = logger.ReturnEntry().Logger

	return &RouterService{
		chi.NewMux(),
	}
}

func (rt *RouterService) MountWith(
	pattern string,
	sub http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) {
	if len(middlewares) > 0 {
		rt.Use(middlewares...)
	}
	rt.Mount(pattern, sub)
}

func MakeSubroute(routes ...HandlerParam) *RouterService {

	router := NewRouterService()

	for _, route := range routes {
		handle := InitController(route.Handler)

		switch route.Method {

		case OPTIONS:
			router.Options(route.Path, handle)
		case GET:
			router.Get(route.Path, handle)
		case HEAD:
			router.Head(route.Path, handle)
		case POST:
			router.Post(route.Path, handle)
		case PUT:
			router.Put(route.Path, handle)
		case PATCH:
			router.Patch(route.Path, handle)
		case DELETE:
			router.Delete(route.Path, handle)
		case TRACE:
			router.Trace(route.Path, handle)

		default:
			router.Get(route.Path, handle)
		}
	}

	return router
}

func InitController(handle handler.APIHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if log == nil {
			log = logger.ReturnEntry().Logger
		}

		op, err := handle(w, r)

		log.Debugf("requested controller: %s", op)

		if err != nil {

			switch e := err.(type) {

			case *handler.APIErrorResponse:
				handler.WriteJSON(w, e.StatusCode, e)
				break

			default:
				errDefault := handler.InternalErrorResponse()
				handler.WriteJSON(w, errDefault.StatusCode, errDefault)
				break
			}

			log.Errorf(
				"API Error - path: %s | controller: %s | error: %s",
				r.URL.Path, op, err.Error(),
			)
		}
	}
}

// func InitAPIRouting(ver int) *APIRouting {
// 	return &APIRouting{
// 		BasePath: fmt.Sprintf("/api/v%d", ver),
// 		MW:       mws.Init(),
// 	}
// }

// func (rt *APIRouting) pathWithBase(pt string) string {
// 	return rt.BasePath + pt
// }

// func (rt *APIRouting) ConfigRoutes() *chi.Mux {

// 	router := rt.setBaseRouting()
// 	protectedRoute := chi.NewRouter()

// 	protectedRoute.Use(rt.MW.AuthorizationJWT, rt.MW.CORSDev)

// 	// controller for auth status checking
// 	protectedRoute.Get("/check", handlers.InitController(frontend.AccessCheck))

// 	protectedRoute.Mount("/apps", setApplicationsRouting())
// 	protectedRoute.Mount("/pve", setProxmoxRouting())
// 	protectedRoute.Mount("/system", setSystemRouting())

// 	router.Mount(rt.BasePath, protectedRoute)

// 	return router
// }

// func (rt *APIRouting) setBaseRouting() *chi.Mux {
// 	chi := chi.NewMux()

// 	chi.Use(rt.MW.PanicRecoverer, rt.MW.Logging, rt.MW.Compressor)

// 	chi.Get("/swagger/*", httpSwagger.WrapHandler)

// 	front := frontend.Init(authorization.InitAuth(configuration.GetConfig()))

// 	chi.Post("/login", handlers.InitController(front.Login))

// 	chi.Get("/", handlers.InitController(front.HTML))
// 	chi.Get("/welcome", handlers.InitController(front.HTML))

// 	chi.Get("/assets/*", handlers.InitController(front.Assets))
// 	chi.Get("/static/*", handlers.InitController(front.Static))
// 	chi.Get("/wallpaper/*", handlers.InitController(front.WallpaperHandle))

// 	return chi
// }

// func setApplicationsRouting() *chi.Mux {

// 	apps, err := appsfile.Init(appsfile.DefaultPath)
// 	if err != nil {
// 		return chi.NewMux()
// 	}

// 	group := applications.Init(apps)

// 	return BuildSubroute(
// 		RoutesConfig{
// 			HandlerParam{"GET", "/table", group.ShowTable},
// 			HandlerParam{"POST", "/table/{topic}", group.AppendApp},
// 			HandlerParam{"DELETE", "/table/{topic}/{number}", group.DeleteApp},
// 		},
// 	)
// }

// func setProxmoxRouting() *chi.Mux {
// 	pve := proxmox.Init(cache.GetEntry())

// 	return BuildSubroute(
// 		RoutesConfig{
// 			HandlerParam{"GET", "/sessions", pve.Sessions},
// 			HandlerParam{"GET", "/{session}/{node}/status", pve.NodeStatus},
// 			HandlerParam{"GET", "/{session}/{node}/devices", pve.DeviceList},

// 			HandlerParam{"GET", "/{session}/{node}/disks", pve.DiskList},
// 			HandlerParam{"GET", "/{session}/{node}/disks/smart", pve.SMART}, // TODO: Fix panic with exists disk

// 			HandlerParam{"GET", "/{session}/{node}/apt/updates", pve.AptUpdates},
// 			HandlerParam{"POST", "/{session}/{node}/apt/update", pve.AptUpdate},

// 			HandlerParam{"POST", "/{session}/{node}/devices/{vmid}/{command}", pve.DeviceCommand},
// 		},
// 	)
// }

// func setSystemRouting() *chi.Mux {
// 	sys := sys.Init(system.NewHostInfoService(), cache.GetEntry())

// 	return BuildSubroute(
// 		RoutesConfig{
// 			HandlerParam{"GET", "/info", sys.HostInfo},

// 			HandlerParam{"GET", "/systemd/status", sys.SystemdUnits},
// 			HandlerParam{"POST", "/systemd/{unit}/{command}", sys.UnitCommand},

// 			HandlerParam{"GET", "/stats", sys.HostStatsWS},
// 			HandlerParam{"GET", "/tty", sys.TtyWS},

// 			HandlerParam{"GET", "/stats/", sys.HostStatsWS},
// 			HandlerParam{"GET", "/tty/", sys.TtyWS},
// 		},
// 	)
// }
