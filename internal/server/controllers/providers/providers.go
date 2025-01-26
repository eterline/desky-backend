package providers

import (
	"net/http"
	"strings"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/server/router/handler"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Exporter interface {
	Get() models.ExportCredentialsStack
	Delete(query int) error
	Append(models.ExporterCredentials) error
}

type ExporterType int

const (
	_ ExporterType = iota
	PVE
	DOCKER
)

type ProvidersControllers struct {
	export map[ExporterType]Exporter
	log    *logrus.Logger
}

func Init() *ProvidersControllers {
	return &ProvidersControllers{
		export: map[ExporterType]Exporter{},
		log:    logger.ReturnEntry().Logger,
	}
}

func (pc *ProvidersControllers) Register(t ExporterType, e Exporter) {
	pc.export[t] = e
}

func (pc *ProvidersControllers) ServiceSessions(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "providers.service-sessions"

	q, err := handler.ParseURLParameters(r, handler.StrOpts("service"))
	if err != nil {
		return op, err
	}

	exp, err := pc.exporter(q.GetStr("service"))
	if err != nil {
		return op, err
	}

	return op, handler.WriteJSON(w, http.StatusOK, exp.Get())
}

func (pc *ProvidersControllers) exporter(name string) (Exporter, error) {
	var sc Exporter = nil

	switch strings.ToLower(name) {
	case "pve":
		sc = pc.export[PVE]
	case "docker":
		sc = pc.export[DOCKER]
	default:
		return nil, handler.NewErrorResponse(
			http.StatusNotFound,
			ErrUnknownService,
		)
	}

	if sc == nil {
		return nil, handler.NewErrorResponse(
			http.StatusNoContent,
			ErrNotConfiguredExport,
		)
	}

	return sc, nil
}
