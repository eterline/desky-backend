package applications

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/apps/appsfile"
)

type AppsService interface {
	Append(topic string, app models.AppDetails) error
	Delete(topic string, topicQuery int) error
	Edit(topic string, topicQuery int, app models.AppDetails) error
	Table() (models.AppsTable, error)
}

type AppsHandlerGroup struct {
	Apps AppsService
}

func Init(service AppsService) *AppsHandlerGroup {
	return &AppsHandlerGroup{
		Apps: service,
	}
}

func (as *AppsHandlerGroup) ShowTable(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.applications.show-table"

	t, err := as.Apps.Table()
	if err == nil {
		handlers.WriteJSON(w, http.StatusOK, t)
	}

	if appsfile.IsAppsFileServiceError(err) {
		err = handlers.NewErrorResponse(
			http.StatusNotImplemented,
			err,
		)
	}

	return op, err
}

func (as *AppsHandlerGroup) AppendApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.applications.append-app"
	data := models.AppDetails{}

	q, err := handlers.ParseURLParameters(r, handlers.StrOpts("topic"))
	if err != nil {
		return op, err
	}

	if err = handlers.DecodeRequest(r, &data); err != nil {
		return op, err
	}

	if err = as.Apps.Append(q.GetStr("topic"), data); err == nil {
		handlers.StatusCreated(w, "app added")
	}

	if appsfile.IsAppsFileServiceError(err) {
		err = handlers.NewErrorResponse(
			http.StatusNotImplemented,
			err,
		)
	}

	return op, err
}

func (as *AppsHandlerGroup) DeleteApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.applications.delete-app"

	q, err := handlers.ParseURLParameters(r, handlers.StrOpts("topic"), handlers.NumOpts("number"))
	if err != nil {
		return op, err
	}

	if err = as.Apps.Delete(q.GetStr("topic"), q.GetInt("number")); err == nil {
		handlers.StatusOK(w, "app deleted")
	}

	if appsfile.IsAppsFileServiceError(err) {
		err = handlers.NewErrorResponse(
			http.StatusNotImplemented,
			err,
		)
	}

	return op, err
}
