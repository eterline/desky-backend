package frontend

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/api/handlers/files"
)

type Handler struct {
	FS       *files.FilesHandler
	HTMLfile string
}

func Init() *Handler {
	return &Handler{
		FS:       files.Init("./web"),
		HTMLfile: "index.html",
	}
}

func (fh *Handler) HTML(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.html"

	fh.FS.ServeFile(w, r, fh.HTMLfile)
	return op, err
}

func (fh *Handler) Static(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.static"

	fh.FS.ServeDir("static").ServeHTTP(w, r)
	return op, err
}

func (fh *Handler) Assets(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.assets"

	fh.FS.ServeDir("assets").ServeHTTP(w, r)
	return op, err
}
