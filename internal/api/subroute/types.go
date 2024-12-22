package subroute

import "net/http"

type HandlerParam struct {
	Method, Path string
	Handler      http.HandlerFunc
}

type RoutesConfig []HandlerParam
