package routing

import (
	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/go-chi/chi"
)

func BuildSubroute(routes RoutesConfig) *chi.Mux {
	chi := chi.NewMux()

	for _, route := range routes {
		h := handlers.InitController(route.Handler)

		switch route.Method {

		case "POST":
			chi.Post(route.Path, h)
			break

		case "DELETE":
			chi.Delete(route.Path, h)
			break

		case "PATCH":
			chi.Patch(route.Path, h)
			break

		default:
			chi.Get(route.Path, h)
			break
		}
	}

	return chi
}
