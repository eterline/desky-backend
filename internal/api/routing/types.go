package routing

import "github.com/eterline/desky-backend/internal/api/handlers"

type HandlerParam struct {
	Method, Path string
	Handler      handlers.APIHandlerFunc
}

type RoutesConfig []HandlerParam
