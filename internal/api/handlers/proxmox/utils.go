package proxmox

import (
	"context"

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
