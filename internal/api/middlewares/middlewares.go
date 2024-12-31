package middlewares

import (
	"net/http"

	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

type MiddleWare struct {
	logger *logrus.Logger
}

func Init() *MiddleWare {
	return &MiddleWare{
		logger: logger.ReturnEntry().Logger,
	}
}

func (mw *MiddleWare) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mw.logger.Infof("[%s] %s - %s", request.Method, request.RemoteAddr, request.RequestURI)
		next.ServeHTTP(writer, request)
	})
}

// TODO:
func (mw *MiddleWare) AuthorizationJWT() {

}
