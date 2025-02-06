package router

import (
	"github.com/eterline/desky-backend/internal/services/router/handler"
	"github.com/go-chi/chi"
)

type RouterService struct {
	*chi.Mux
}

type HandlerMethod int

type HandlerParam struct {
	Method  HandlerMethod
	Path    string
	Handler handler.APIHandlerFunc
}

func NewHandler(m HandlerMethod, p string, h handler.APIHandlerFunc) HandlerParam {
	return HandlerParam{
		Method:  m,
		Path:    p,
		Handler: h,
	}
}
