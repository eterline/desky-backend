package files

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type FilesHandlerGroup struct {
	BasePath string
}

func Init(base string) *FilesHandlerGroup {
	log = logger.ReturnEntry().Logger

	return &FilesHandlerGroup{
		BasePath: base,
	}
}

func (fh *FilesHandlerGroup) PathWithBase(path string) string {
	return filepath.Join(fh.BasePath, path)
}

func (fh *FilesHandlerGroup) ServeDir(dir string) http.Handler {
	op := "handlers.files.servedir"
	log.Debugf("requested controller: %s", op)

	fs := http.FileServer(http.Dir(fh.PathWithBase(dir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" || r.URL.Path[0] == '.' {

			err := handlers.ForbiddenRequestResponse()
			handlers.WriteJSON(w, err.StatusCode, err.Message)

		}
		http.StripPrefix(fmt.Sprintf("/%s/", dir), fs).ServeHTTP(w, r)
	})
}

func (fh *FilesHandlerGroup) ServeFile(w http.ResponseWriter, r *http.Request, filePath string) {
	op := "handlers.files.servefile"
	log.Debugf("requested controller: %s", op)

	http.ServeFile(w, r, fh.PathWithBase(filePath))
}
