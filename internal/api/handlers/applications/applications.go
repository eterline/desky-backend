package applications

import (
	"net/http"
	"path/filepath"

	"github.com/eterline/desky-backend/internal/api/handlers"
)

type AppsService struct {
	filePath string
}

func Init(path string) *AppsService {
	return &AppsService{
		filePath: filepath.Join(path),
	}
}

func (as *AppsService) ReturnAppsTable(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.applications.returnappstable"

	t, err := as.getAppTable()
	if err == nil {
		handlers.WriteJSON(w, http.StatusOK, t)
	}

	return op, err
}
