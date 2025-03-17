package controllers

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/apps/appsfile"
	"github.com/eterline/desky-backend/internal/services/handler"
)

type AppsService interface {
	Append(topic string, app models.AppDetails) error
	Delete(topic string, topicQuery int) error
	Edit(app *models.AppDetails) error
	Table() (models.AppsTable, error)
}

type AppsDeleter interface {
	DeleteApp(id uint) error
}

type AppsHandlerGroup struct {
	Apps AppsService
	Del  AppsDeleter
}

func InitApplications(service AppsService, del AppsDeleter) *AppsHandlerGroup {
	return &AppsHandlerGroup{
		Apps: service,
		Del:  del,
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
func (as *AppsHandlerGroup) CreateApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.applications.create-app"

	q, err := handler.ParseURLParameters(r, handler.StrOpts("topic"))
	if err != nil {
		return op, err
	}

	data := new(models.AppDetails)
	handler.DecodeRequest(r, data)

	if err := handler.Validate(data); err != nil {
		return op, err
	}

	if err = as.Apps.Append(q.GetStr("topic"), *data); err == nil {
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

func (as *AppsHandlerGroup) EditApp(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.applications.edit-app"

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	form := new(models.AppUpdateFrom)
	handler.DecodeRequest(r, form)

	if err := handler.Validate(form); err != nil {
		return op, err
	}

	data := &models.AppDetails{
		ID:          uint(q.GetInt("id")),
		Name:        form.Name,
		Description: form.Description,
		Link:        form.Link,
		Icon:        form.Icon,
	}

	if err = as.Apps.Edit(data); err == nil {
		handler.StatusCreated(w, "app edited")
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
//	@Router			/apps/table/{id} [delete]
func (as *AppsHandlerGroup) DeleteAppById(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.applications.delete-app-by-id"

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	id := uint(q.GetInt("id"))

	if err = as.Del.DeleteApp(id); err != nil {
		return op, err
	}

	return op, handler.StatusOK(w, "app deleted")
}
