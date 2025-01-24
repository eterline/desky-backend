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

// Sessions godoc
//
//	@Summary		Sessions
//	@Description	Proxmox sessions
//	@Tags			pve
//
//	@Accept			json
//	@Produce		json
//	@Failure		500	{object}	handlers.APIErrorResponse
//	@Success		200	{array}		[]ve.SessionNodes
//	@Router			/pve/sessions [get]
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

// NodeStatus godoc
//
//	@Summary		NodeStatus
//	@Description	Retrieve detailed status information for a Proxmox VE node, including load, filesystem, RAM, CPU, and uptime.
//	@Tags			pve
//
//	@Produce		json
//	@Param			session	path		string	true	"Session ID"
//	@Param			node	path		string	true	"Node Name"
//	@Success		200		{object}	models.PVENodeStatus
//	@Failure		400		{object}	handlers.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handlers.APIErrorResponse	"Internal server error"
//	@Router			/pve/{session}/{node}/status [get]
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

// DeviceList godoc
//
//	@Tags			pve
//	@Summary		DeviceList
//	@Description	Getting ve devices list information.
//
//	@Produce		json
//	@Param			session	path		string	true	"Session ID"
//	@Param			node	path		string	true	"Node Name"
//	@Success		200		{object}	models.DevicesList
//	@Failure		400		{object}	handlers.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handlers.APIErrorResponse	"Internal server error"
//	@Router			/pve/{session}/{node}/devices [get]
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

// DeviceCommand godoc
//
//	@Tags			pve
//	@Summary		DeviceCommand
//	@Description	execute device command
//
//	@Produce		json
//	@Param			session	path		string	true	"Session ID"
//	@Param			node	path		string	true	"Node Name"
//	@Param			vmid	path		string	true	"VMID"
//	@Param			command	query		string	true	"VM Command"	Enums(stop, start, shutdown, suspend, resume)
//	@Success		200		{object}	handlers.APIResponse
//	@Failure		400		{object}	handlers.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handlers.APIErrorResponse	"Internal server error"
//	@Failure		501		{object}	handlers.APIErrorResponse	"Uninplemented"
//	@Router			/pve/{session}/{node}/devices/{vmid}/{command} [post]
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

// AptUpdates godoc
//
//	@Tags			pve
//	@Summary		AptUpdates
//	@Description	getting apt proxmox update list
//
//	@Produce		json
//	@Param			session	path		string	true	"Session ID"
//	@Param			node	path		string	true	"Node Name"
//	@Success		200		{object}	nodes.AptUpdates
//	@Failure		400		{object}	handlers.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handlers.APIErrorResponse	"Internal server error"
//	@Router			/pve/{session}/{node}/apt/updates [get]
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

// AptUpdate godoc
//
//	@Tags			pve
//	@Summary		AptUpdate
//	@Description	update proxmox apt packages
//
//	@Produce		json
//	@Param			session	path		string	true	"Session ID"
//	@Param			node	path		string	true	"Node Name"
//	@Success		200		{object}	handlers.APIResponse
//	@Failure		400		{object}	handlers.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handlers.APIErrorResponse	"Internal server error"
//	@Router			/pve/{session}/{node}/apt/update [post]
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

// DiskList godoc
//
//	@Tags			pve
//	@Summary		DiskList
//	@Description	getting disk list
//
//	@Produce		json
//	@Param			session	path		string	true	"Session ID"
//	@Param			node	path		string	true	"Node Name"
//	@Success		200		{object}	nodes.DisksInfo
//	@Failure		400		{object}	handlers.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handlers.APIErrorResponse	"Internal server error"
//	@Router			/pve/{session}/{node}/disks [get]
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

// smart godoc
//
//	@Tags			pve
//	@Summary		smart
//	@Description	getting disk SMART info
//
//	@Produce		json
//	@Param			session	path		string	true	"Session ID"
//	@Param			node	path		string	true	"Node Name"
//	@Param			dev		query		string	true	"device path"
//	@Success		200		{object}	nodes.Smart
//	@Failure		400		{object}	handlers.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handlers.APIErrorResponse	"Internal server error"
//	@Router			/pve/{session}/{node}/disks/smart [get]
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
