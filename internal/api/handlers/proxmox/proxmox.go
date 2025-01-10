package proxmox

import (
	"context"
	"net/http"
	"time"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/services/ve"
	"github.com/eterline/desky-backend/pkg/logger"
)

type Cacher interface {
	GetValue(key any) any
	PushValue(key, value any)
	OlderThanAndExists(key any, duration time.Duration) bool
}

type ProxmoxHandlerGroup struct {
	Provider *ve.VEService
	Cache    Cacher
}

func Init(ch Cacher) *ProxmoxHandlerGroup {
	pve := ve.Init()

	logger.ReturnEntry().Infof("%v pve sessions connected with %v errors", pve.AvailSessions(), pve.ErrCount)

	return &ProxmoxHandlerGroup{
		Provider: pve,
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

func (ph *ProxmoxHandlerGroup) Sessions(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.sessions"

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	list, err := ph.Provider.SessionList()
	if err == nil {
		err = handlers.WriteJSON(w, http.StatusOK, list)
	}

	return op, err
}

func (ph *ProxmoxHandlerGroup) NodeStatus(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.node-status"

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	q, err := handlers.ParseURLParameters(r, handlers.NumOpts("session"), handlers.StrOpts("node"))
	if err != nil {
		return op, err
	}

	pveSession, err := ph.Provider.GetSession(q.GetInt("session"))
	if err != nil {
		return op, err
	}

	status, err := pveSession.NodeStatus(context.Background(), q.GetStr("node"))
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

	q, err := handlers.ParseURLParameters(r, handlers.NumOpts("session"), handlers.StrOpts("node"))
	if err != nil {
		return op, err
	}

	pveSession, err := ph.Provider.GetSession(q.GetInt("session"))
	if err != nil {
		return op, err
	}

	devices, err := pveSession.DeviceList(context.Background(), q.GetStr("node"))
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

	q, err := handlers.ParseURLParameters(r,
		handlers.NumOpts("session", "vmid"),
		handlers.StrOpts("command", "node"),
	)
	if err != nil {
		return op, err
	}

	pveSession, err := ph.Provider.GetSession(q.GetInt("session"))
	if err != nil {
		return op, err
	}

	dev, err := pveSession.ResolveDevice(q.GetStr("node"), q.GetInt("vmid"))
	if err != nil {
		return op, err
	}

	err, found := execDeviceCommand(dev, q.GetStr("command"), context.Background())

	if err != nil {

		if found {
			return op, handlers.NewErrorResponse(
				http.StatusNotImplemented,
				ErrActionCannotComplete(q.GetInt("vmid")),
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

func (ph *ProxmoxHandlerGroup) AptUpdates(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.apt-updates"

	ctx := context.Background()

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	node, err := ph.parseVeNode(r)
	if err != nil {
		return op, err
	}

	u, err := node.GetAptUpdates(ctx)
	if err != nil {
		return op, err
	}

	return op, handlers.WriteJSON(w, http.StatusOK, u.Data)
}

func (ph *ProxmoxHandlerGroup) AptUpdate(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.apt-updates"
	ctx := context.Background()

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	node, err := ph.parseVeNode(r)
	if err != nil {
		return op, err
	}

	if _, err := node.AptUpgrade(ctx); err != nil {
		return op, err
	}

	return op, handlers.StatusOK(w, "update task successfully")
}

func (ph *ProxmoxHandlerGroup) DiskList(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.disk-list"
	ctx := context.Background()

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	node, err := ph.parseVeNode(r)
	if err != nil {
		return op, err
	}

	lst, err := node.Disks(ctx)
	if err != nil {
		return op, err
	}

	return op, handlers.WriteJSON(w, http.StatusOK, lst.Data)
}

func (ph *ProxmoxHandlerGroup) SMART(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.proxmox.disk-list"
	ctx := context.Background()

	r.ParseForm()

	if err := ph.invalidSessionsHandler(w); err != nil {
		return op, err
	}

	node, err := ph.parseVeNode(r)
	if err != nil {
		return op, err
	}

	disk, err := node.DiskByDevPath(ctx, r.FormValue("dev"))
	if err != nil {
		return op, err
	}

	smart, err := disk.SMART(ctx)
	if err != nil {
		return op, err
	}

	return op, handlers.WriteJSON(w, http.StatusOK, smart)
}
