package routing

import (
	"fmt"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/api/handlers/applications"
	"github.com/eterline/desky-backend/internal/api/handlers/frontend"
	"github.com/eterline/desky-backend/internal/api/middlewares"
	"github.com/go-chi/chi"
)

type APIRouting struct {
	BasePath string
}

func InitAPIRouting(ver int) *APIRouting {
	return &APIRouting{
		BasePath: fmt.Sprint("/api/v%w", ver),
	}
}

func (rt *APIRouting) pathWithBase(pt string) string {
	return rt.BasePath + pt
}

func (rt *APIRouting) ConfigRoutes() *chi.Mux {
	router := SetBaseRouting()

	router.Mount("/apps", SetApplicationsRouting())

	return router
}

func SetBaseRouting() *chi.Mux {
	chi := chi.NewMux()
	mw := middlewares.Init()

	chi.Use(mw.Logging)

	front := frontend.Init()

	chi.Get("/", handlers.InitController(front.HTML))
	chi.Get("/assets/*", handlers.InitController(front.Assets))
	chi.Get("/static/*", handlers.InitController(front.Static))

	return chi
}

func SetApplicationsRouting() *chi.Mux {
	as := applications.Init("apps.json")

	return BuildSubroute(
		RoutesConfig{
			HandlerParam{"GET", "/table", as.ReturnAppsTable},
		},
	)
}
