package middlewares

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/services/router/handler"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type MiddlewareService struct {
	logger *logrus.Logger
	config configuration.Configuration
}

func Init(c configuration.Configuration) *MiddlewareService {
	return &MiddlewareService{
		logger: logger.ReturnEntry().Logger,
		config: c,
	}
}

func (mw *MiddlewareService) CORSDev(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		next.ServeHTTP(writer, request)

		if mw.config.DevelopEnv {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		}
	})
}

func (mw *MiddlewareService) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		mw.logger.Infof("Request [HTTP] [%s] %s - %s", request.Method, request.RemoteAddr, request.RequestURI)

		if websocket.IsWebSocketUpgrade(request) {

			mw.logger.Infof("Request [WS] [%s] %s - %s", request.Method, request.RemoteAddr, request.RequestURI)
			next.ServeHTTP(writer, request)

			return
		}

		rw := NewResponseWriter(writer)
		next.ServeHTTP(rw, request)

		mw.logger.Infof("Response [HTTP] [%d] - %s - %s", rw.statusCode, request.RemoteAddr, request.RequestURI)

	})
}

func (mw *MiddlewareService) Compressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if websocket.IsWebSocketUpgrade(request) {

			next.ServeHTTP(writer, request)
			return
		}

		middleware.Compress(5, "text/html", "text/css")
		(next).ServeHTTP(writer, request)
	})
}

func (mw *MiddlewareService) PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				mw.logger.Infof("handler panic recovered: in 500 status code: %v", r)

				e := handler.InternalErrorResponse()
				handler.WriteJSON(writer, e.StatusCode, e)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}
