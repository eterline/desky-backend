package server

import (
	"crypto/tls"
	"net/http"

	"github.com/eterline/desky-backend/internal/server/router"
)

type Server struct {
	srv  *http.Server
	cert string
	key  string
}

// New - create new server object
func New(
	socket string,
	cert string, key string,
	serverName string,
) *Server {
	return &Server{
		srv: &http.Server{
			Addr: socket,
			TLSConfig: &tls.Config{
				ServerName: serverName,
			},
		},
		cert: cert,
		key:  key,
	}
}

// Router - set up server Handler router
func (s *Server) Router(r http.Handler) {
	s.srv.Handler = r
}

func (s *Server) Run(ssl bool) error {

	var err error

	if s.srv.Handler == nil {
		panic("couldn't start server without handler")
	}

	if ssl {
		err = s.srv.ListenAndServeTLS(s.cert, s.key)
	} else {
		err = s.srv.ListenAndServe()
	}

	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop - stop the http server
func (s *Server) Stop() error {
	return s.srv.Close()
}

func ConfigRoutes() (r *router.RouterService) {
	r = router.NewRouterService()

	public(r)

	r.MountWith("/api", api())

	return
}
