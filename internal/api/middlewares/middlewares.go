package middlewares

import (
	"net/http"
	"strings"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/services/authorization"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type JWTValidator interface {
	TokenIsValid(string) bool
}

type MiddleWare struct {
	logger      *logrus.Logger
	jwt         JWTValidator
	authEnabled bool
	isDev       bool
}

func Init() *MiddleWare {

	config := configuration.GetConfig()

	return &MiddleWare{
		logger:      logger.ReturnEntry().Logger,
		jwt:         authorization.NewJWTauth(config.Server.JWTSecretBytes()),
		authEnabled: config.Auth.Enabled,
		isDev:       config.DevMode,
	}
}

func (mw *MiddleWare) CORSDev(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		next.ServeHTTP(writer, request)

		if mw.isDev {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		}
	})
}

func (mw *MiddleWare) Logging(next http.Handler) http.Handler {
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

func (mw *MiddleWare) Compressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if websocket.IsWebSocketUpgrade(request) {

			next.ServeHTTP(writer, request)
			return
		}

		middleware.Compress(5, "text/html", "text/css")
		(next).ServeHTTP(writer, request)
	})
}

func (mw *MiddleWare) AuthorizationJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if !mw.authEnabled {
			next.ServeHTTP(writer, request)
			return
		}

		token := parseBearer(request, "DeskyJWT")

		if !mw.jwt.TokenIsValid(token) {

			mw.logger.Infof("[JWT] %s - unauthorized api request - %s | header token: %s", request.RemoteAddr, request.RequestURI, token)

			e := handlers.NewErrorResponse(http.StatusUnauthorized, ErrNotAuthorized)
			handlers.WriteJSON(writer, e.StatusCode, e)

			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (mw *MiddleWare) PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		next.ServeHTTP(writer, request)

		if r := recover(); r != nil {
			mw.logger.Infof("handler panic recovered: in 501 status code | description: %v", r)

			e := handlers.InternalErrorResponse()
			handlers.WriteJSON(writer, e.StatusCode, e)
		}
	})
}

func parseBearer(r *http.Request, key string) string {

	b := r.Header.Get(key)

	if b == "" {
		r.ParseForm()
		b = r.FormValue(key)
	}

	return strings.ReplaceAll(b, "Bearer ", "")
}
