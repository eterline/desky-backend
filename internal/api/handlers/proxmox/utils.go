package proxmox

import (
	"context"
	"net/http"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/nodes"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/virtual"
)

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

func (ph *ProxmoxHandlerGroup) parseVeNode(r *http.Request) (*nodes.ProxmoxNode, error) {

	q, err := handlers.ParseURLParameters(
		r, handlers.NumOpts("session"),
		handlers.StrOpts("node"),
	)

	if err != nil {
		return nil, err
	}

	pveSession, err := ph.Provider.GetSession(q.GetInt("session"))
	if err != nil {
		return nil, err
	}

	node, err := pveSession.ResolveNode(q.GetStr("node"))
	if err != nil {
		return nil, err
	}

	return node, nil
}
