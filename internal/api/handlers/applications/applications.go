package applications

import (
	"net/http"
	"path/filepath"
	"sync"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/go-chi/chi"
)

type AppsHandlerGroup struct {
	filePath string
	mutx     sync.Mutex
}

func Init(path string) *AppsHandlerGroup {
	return &AppsHandlerGroup{
		filePath: filepath.Join(path),
	}
}

func (as *AppsHandlerGroup) ReturnAppsTable(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.applications.returnappstable"

	t, err := as.getAppTable()
	if err == nil {
		handlers.WriteJSON(w, http.StatusOK, t)
	}

	return op, err
}

func (as *AppsHandlerGroup) AppendApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.applications.appendapp"

	data := AppDetails{}
	topicName := chi.URLParam(r, "topic")

	if err = handlers.DecodeRequest(r, &data); err != nil {
		return op, err
	}

	if err = as.addApp(topicName, data); err == nil {
		handlers.StatusCreated(w, "app added")
	}

	return op, err
}

func (as *AppsHandlerGroup) DeleteApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.applications.deleteapp"

	topicName := chi.URLParam(r, "topic")
	appNumber := chi.URLParam(r, "number")

	if err = as.rmApp(topicName, appNumber); err == nil {
		handlers.StatusOK(w, "app deleted")
	}

	return op, err
}
