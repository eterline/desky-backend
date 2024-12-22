package subroute

import "github.com/go-chi/chi"

func ReturnConfiguredSubroute(routes RoutesConfig) *chi.Mux {
	chi := chi.NewMux()

	for _, route := range routes {
		switch route.Method {

		case "POST":
			chi.Post(route.Path, route.Handler)
			break

		case "DELETE":
			chi.Delete(route.Path, route.Handler)
			break

		case "PATCH":
			chi.Patch(route.Path, route.Handler)
			break

		default:
			chi.Get(route.Path, route.Handler)
			break
		}
	}

	return chi
}
