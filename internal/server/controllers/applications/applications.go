package applications

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/server/router/handler"
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

// ShowTable godoc
//
//	@Summary		ShowTable
//	@Description	Showing apps table with their info
//	@Tags			applications
//
//	@Accept			json
//	@Produce		json
//	@Failure		501	{object}	handler.APIErrorResponse
//	@Success		200	{object}	models.AppsTable
//	@Router			/apps/table [get]
func (as *AppsHandlerGroup) ShowTable(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.applications.show-table"

	t, err := as.Apps.Table()
	if err == nil {
		handler.WriteJSON(w, http.StatusOK, t)
	}

	if appsfile.IsAppsFileServiceError(err) {
		err = handler.NewErrorResponse(
			http.StatusNotImplemented,
			err,
		)
	}

	return op, err
}

// AppendApp godoc
//
//	@Summary		AppendApp
//	@Description	Adding app
//	@Tags			applications
//
//	@Param			request	body	models.AppDetails	true	"app params"
//	@Accept			json
//	@Produce		json
//	@Failure		501	{object}	handler.APIErrorResponse
//	@Success		200	{object}	handler.APIResponse
//	@Router			/apps/table/{topic} [post]
func (as *AppsHandlerGroup) AppendApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.applications.append-app"
	data := models.AppDetails{}

	q, err := handler.ParseURLParameters(r, handler.StrOpts("topic"))
	if err != nil {
		return op, err
	}

	if err = handler.DecodeRequest(r, &data); err != nil {
		return op, err
	}

	if err = as.Apps.Append(q.GetStr("topic"), data); err == nil {
		handler.StatusCreated(w, "app added")
	}

	if appsfile.IsAppsFileServiceError(err) {
		err = handler.NewErrorResponse(
			http.StatusNotImplemented,
			err,
		)
	}

	return op, err
}

// DeleteApp godoc
//
//	@Summary		DeleteApp
//	@Description	Deleting app
//	@Tags			applications
//
//	@Accept			json
//	@Produce		json
//	@Failure		501	{object}	handler.APIErrorResponse
//	@Success		200	{object}	handler.APIResponse
//
//	@Router			/apps/table/{topic}/{number} [delete]
func (as *AppsHandlerGroup) DeleteApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.applications.delete-app"

	q, err := handler.ParseURLParameters(r, handler.StrOpts("topic"), handler.NumOpts("number"))
	if err != nil {
		return op, err
	}

	if err = as.Apps.Delete(q.GetStr("topic"), q.GetInt("number")); err == nil {
		handler.StatusOK(w, "app deleted")
	}

	if appsfile.IsAppsFileServiceError(err) {
		err = handler.NewErrorResponse(
			http.StatusNotImplemented,
			err,
		)
	}

	return op, err
}
