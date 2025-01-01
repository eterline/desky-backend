package proxmox

import (
	"context"
	"net/http"
	"time"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/services/ve"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/virtual"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type Cacher interface {
	GetValue(key any) any
	PushValue(key, value any)
	OlderThanAndExists(key any, duration time.Duration) bool
}

type ProxmoxHandlerGroup struct {
	Provider *ve.ProxmoxProvider
	Cache    Cacher
}

func Init(ch Cacher) *ProxmoxHandlerGroup {
	log = logger.ReturnEntry().Logger

	provider, err := ve.NewProvide()
	if err != nil {
		log.Errorf("proxmox provider initialization error: %s", err.Error())
	}
	log.Infof("%v pve sessions connected with %v errors", provider.AvailSessions(), provider.ErrCount)

	return &ProxmoxHandlerGroup{
		Provider: provider,
		Cache:    ch,
	}
}

func (ph *ProxmoxHandlerGroup) invalidSessionsHandler(w http.ResponseWriter) error {
	if !ph.Provider.AnyValidConns() {
		return handlers.NewErrorResponse(
			http.StatusLocked,
			ErrNoSessions,
		)
	}
	return nil
}

func (ph *ProxmoxHandlerGroup) NodeStatus(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.node-status"

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	q, err := handlers.QueryURLNumeredParameters(r, "session")
	if err != nil {
		return op, err
	}

	qStr, err := handlers.QueryURLParameters(r, "node")
	if err != nil {
		return op, err
	}

	pveSession, err := ph.Provider.GetSession(q["session"])
	if err != nil {
		return op, err
	}

	status, err := pveSession.NodeStatus(context.Background(), qStr["node"])
	if err != nil {
		return op, err
	}

	err = handlers.WriteJSON(w, http.StatusOK, status)

	return op, err
}

func (ph *ProxmoxHandlerGroup) DeviceList(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-list"

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	q, err := handlers.QueryURLNumeredParameters(r, "session")
	if err != nil {
		return op, err
	}

	qStr, err := handlers.QueryURLParameters(r, "node")
	if err != nil {
		return op, err
	}

	pveSession, err := ph.Provider.GetSession(q["session"])
	if err != nil {
		return op, err
	}

	devices, err := pveSession.DeviceList(context.Background(), qStr["node"])
	if err != nil {
		return op, err
	}

	err = handlers.WriteJSON(w, http.StatusOK, devices)

	return op, err
}

func (ph *ProxmoxHandlerGroup) DeviceCommand(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-command"

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	qStr, err := handlers.QueryURLParameters(r, "command")
	if err != nil {
		return op, err
	}

	q, err := handlers.QueryURLNumeredParameters(r, "session", "vmid")
	if err != nil {
		return op, err
	}

	pveSession, err := ph.Provider.GetSession(q["session"])
	if err != nil {
		return op, err
	}

	dev, err := pveSession.ResolveDevice(chi.URLParam(r, "node"), q["vmid"])
	if err != nil {
		return op, err
	}

	err, found := execDeviceCommand(dev, qStr["command"], context.Background())

	if err != nil {
		log.Errorf("proxmox error: %s", err.Error())

		if found {
			return op, handlers.NewErrorResponse(
				http.StatusNotImplemented,
				ErrActionCannotComplete(q["vmid"]),
			)
		} else {
			return op, handlers.NewErrorResponse(
				http.StatusBadRequest,
				err,
			)
		}

	}

	return op, handlers.StatusOK(w, "operation successfully")
}

func execDeviceCommand(dev *virtual.VirtMachine, command string, ctx context.Context) (err error, foundCommand bool) {

	foundCommand = false

	switch command {

	case "start":
		err = dev.Start(ctx)
		break

	case "shutdown":
		err = dev.Shutdown(ctx)
		break

	case "stop":
		err = dev.Stop(ctx)
		break

	case "suspend":
		err = dev.Suspend(ctx)
		break

	case "resume":
		err = dev.Resume(ctx)
		break

	default:
		return ErrUnknownCommand, false
	}

	return err, true
}
