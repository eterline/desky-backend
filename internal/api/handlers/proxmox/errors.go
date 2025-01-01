package proxmox

import (
	"errors"
	"fmt"
)

var (
	ErrActionCannotComplete = func(vmid int) error {
		return fmt.Errorf("VMID %v: action can not be complete", vmid)
	}
	ErrUnknownCommand = errors.New("unknown device command")

	ErrNoSessions = errors.New("no valid sessions")
)
