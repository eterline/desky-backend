package files

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type FilesHandler struct {
	BasePath string
}

func Init(base string) *FilesHandler {
	log = logger.ReturnEntry().Logger

	return &FilesHandler{
		BasePath: base,
	}
}

func (fh *FilesHandler) PathWithBase(path string) string {
	return filepath.Join(fh.BasePath, path)
}

func (fh *FilesHandler) ServeDir(dir string) http.Handler {
	op := "handlers.files.servedir"
	log.Debugf("requested controller: %s", op)

	fs := http.FileServer(http.Dir(fh.PathWithBase(dir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path[0] == '.' {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.StripPrefix(fmt.Sprintf("/%s/", dir), fs).ServeHTTP(w, r)
	})
}

func (fh *FilesHandler) ServeFile(w http.ResponseWriter, r *http.Request, filePath string) {
	op := "handlers.files.servefile"
	log.Debugf("requested controller: %s", op)

	http.ServeFile(w, r, fh.PathWithBase(filePath))
}
