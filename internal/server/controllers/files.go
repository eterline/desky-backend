package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/eterline/desky-backend/internal/services/handler"
)

type FilesHandlerGroup struct {
	BasePath string
}

func InitFiles(base string) *FilesHandlerGroup {
	return &FilesHandlerGroup{
		BasePath: base,
	}
}

func (fh *FilesHandlerGroup) PathWithBase(path string) string {
	return filepath.Join(fh.BasePath, path)
}

func (fh *FilesHandlerGroup) ServeDir(dir string) http.Handler {

	fs := http.FileServer(http.Dir(fh.PathWithBase(dir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" || r.URL.Path[0] == '.' || strings.HasPrefix(r.URL.Path, "/.") {

			err := handler.ForbiddenRequestResponse()
			handler.WriteJSON(w, err.StatusCode, err.Message)
			return
		}
		http.StripPrefix(fmt.Sprintf("/%s/", dir), fs).ServeHTTP(w, r)
	})
}

func (fh *FilesHandlerGroup) ServeFile(w http.ResponseWriter, r *http.Request, filePath string) {
	http.ServeFile(w, r, fh.PathWithBase(filePath))
}
