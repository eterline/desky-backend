package frontend

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/api/handlers/files"
)

type FrontendHandlerGroup struct {
	FS       *files.FilesHandlerGroup
	HTMLfile string
}

func Init() *FrontendHandlerGroup {
	return &FrontendHandlerGroup{
		FS:       files.Init("./web"),
		HTMLfile: "index.html",
	}
}

func (fh *FrontendHandlerGroup) HTML(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.html"

	fh.FS.ServeFile(w, r, fh.HTMLfile)
	return op, err
}

func (fh *FrontendHandlerGroup) Static(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.static"

	fh.FS.ServeDir("static").ServeHTTP(w, r)
	return op, err
}

func (fh *FrontendHandlerGroup) Assets(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.assets"

	fh.FS.ServeDir("assets").ServeHTTP(w, r)
	return op, err
}
