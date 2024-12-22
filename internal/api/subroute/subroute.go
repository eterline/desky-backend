package subroute

import (
	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/go-chi/chi"
)

var handlerProvider *handlers.ApiHandler

func ConfigAPI(h *handlers.ApiHandler) *chi.Mux {
	handlerProvider = h

	chi := chi.NewMux()

	chi.Mount("/apps", appsSubroute())
	chi.Mount("/proxmox", appsSubroute())
	chi.Mount("/docker", appsSubroute())
	chi.Mount("/host", appsSubroute())

	return chi
}

func appsSubroute() *chi.Mux {
	return ReturnConfiguredSubroute(
		RoutesConfig{
			HandlerParam{"GET", "/list", handlerProvider.AppsList},
			HandlerParam{"POST", "/app", handlerProvider.CreateInAppsList},
			HandlerParam{"DELETE", "/app", handlerProvider.DeleteFromAppsList},
		},
	)
}
