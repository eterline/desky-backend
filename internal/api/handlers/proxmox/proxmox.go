package proxmox

import (
	"context"
	"net/http"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/services/ve"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type ProxmoxHandlerGroup struct {
	Provider *ve.ProxmoxProvider
}

func Init() (*ProxmoxHandlerGroup, error) {
	log = logger.ReturnEntry().Logger

	provider, err := ve.NewProvide()
	if err != nil {
		return nil, err
	}

	return &ProxmoxHandlerGroup{
		Provider: provider,
	}, nil
}

func (ph *ProxmoxHandlerGroup) NodeList(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.node-list"

	return op, err
}

func (ph *ProxmoxHandlerGroup) NodeStatus(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.node-status"

	return op, err
}

func (ph *ProxmoxHandlerGroup) DeviceList(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-list"

	return op, err
}

func (ph *ProxmoxHandlerGroup) DeviceStart(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-start"

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

	if err = dev.Start(context.Background()); err != nil {

		log.Errorf("proxmox error: %s", err.Error())

		return op, handlers.NewErrorResponse(
			http.StatusNotImplemented,
			ErrActionCannotComplete(q["vmid"]),
		)
	}

	return op, nil
}

func (ph *ProxmoxHandlerGroup) DeviceShutdown(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-shutdown"

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

	if err = dev.Shutdown(context.Background()); err != nil {

		log.Errorf("proxmox error: %s", err.Error())

		return op, handlers.NewErrorResponse(
			http.StatusNotImplemented,
			ErrActionCannotComplete(q["vmid"]),
		)
	}

	return op, nil
}

func (ph *ProxmoxHandlerGroup) DeviceStop(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-stop"

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

	if err = dev.Stop(context.Background()); err != nil {

		log.Errorf("proxmox error: %s", err.Error())

		return op, handlers.NewErrorResponse(
			http.StatusNotImplemented,
			ErrActionCannotComplete(q["vmid"]),
		)
	}

	return op, nil
}

func (ph *ProxmoxHandlerGroup) DeviceSuspend(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-suspend"

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

	if err = dev.Suspend(context.Background()); err != nil {

		log.Errorf("proxmox error: %s", err.Error())

		return op, handlers.NewErrorResponse(
			http.StatusNotImplemented,
			ErrActionCannotComplete(q["vmid"]),
		)
	}

	return op, nil
}

func (ph *ProxmoxHandlerGroup) DeviceResume(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.device-resume"

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

	if err = dev.Resume(context.Background()); err != nil {

		log.Errorf("proxmox error: %s", err.Error())

		return op, handlers.NewErrorResponse(
			http.StatusNotImplemented,
			ErrActionCannotComplete(q["vmid"]),
		)
	}

	return op, nil
}
