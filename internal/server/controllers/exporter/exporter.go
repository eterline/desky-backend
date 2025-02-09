package exporter

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/router/handler"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = nil

type Exporter interface {
	Append(form models.ExporterForm) error
	Delete(id int) error
	Services(exporter models.ExporterTypeString) ([]models.ExporterInfo, error)
}

type ExporterControllers struct {
	service Exporter
}

func Init(ex Exporter) *ExporterControllers {

	log = logger.ReturnEntry().Logger

	return &ExporterControllers{
		service: ex,
	}
}

func (exc *ExporterControllers) ListAll(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "exporter.list-all"

	exports, err := exc.service.Services("")
	if err != nil {
		log.Error(err)
		return op, err
	}

	return op, handler.WriteJSON(w, http.StatusOK, exports)
}

func (exc *ExporterControllers) ListType(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "exporter.list-type"

	q, err := handler.ParseURLParameters(r, handler.StrOpts("service"))
	if err != nil {
		return op, err
	}

	exports, err := exc.service.Services(models.ExporterTypeString(q.GetStr("service")))
	if err != nil {
		log.Error(err)
		return op, err
	}

	return op, handler.WriteJSON(w, http.StatusOK, exports)
}

func (exc *ExporterControllers) Append(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "exporter.append"

	q, err := handler.ParseURLParameters(r, handler.StrOpts("service"))
	if err != nil {
		return op, err
	}

	if err := exc.appendExporter(r, models.ExporterTypeString(q.GetStr("service"))); err != nil {
		return op, err
	}

	return op, handler.StatusOK(w, "exporter added")
}

func (exc *ExporterControllers) Delete(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "exporter.delete"

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	if err := exc.service.Delete(q.GetInt("id")); err != nil {
		return op, err
	}

	return op, handler.StatusOK(w, "exporter deleted")
}

func (exc *ExporterControllers) appendExporter(r *http.Request, typeStr models.ExporterTypeString) error {

	switch typeStr {

	case models.ExporterProxmoxType:
		form := new(models.ProxmoxFormExport)
		handler.DecodeRequest(r, form)
		return exc.validForm(form)

	case models.ExporterDockerType:
		form := new(models.DockerFormExport)
		handler.DecodeRequest(r, form)
		return exc.validForm(form)

	default:
		return handler.NotFoundPageResponse()
	}
}

func (exc *ExporterControllers) validForm(form models.ExporterForm) error {
	if err := handler.Validate(form); err != nil {
		return err
	}
	if err := exc.service.Append(form); err != nil {
		return err
	}
	return nil
}
