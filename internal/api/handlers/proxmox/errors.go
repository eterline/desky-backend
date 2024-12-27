package proxmox

import "fmt"

var (
	ErrActionCannotComplete = func(vmid int) error {
		return fmt.Errorf("VMID %v: action can not be complete", vmid)
	}
)
