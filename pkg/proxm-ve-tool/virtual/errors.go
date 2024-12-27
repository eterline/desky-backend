package virtual

import "fmt"

var (
	ErrVirtualNotExists = func(vmid int) error {
		return fmt.Errorf("virtual device with VMID: %v does not exists", vmid)
	}

	ErrNotQEMU = func(vmid int) error {
		return fmt.Errorf("virtual device with VMID: %v does not qemu type", vmid)
	}

	ErrNotLXC = func(vmid int) error {
		return fmt.Errorf("virtual device with VMID: %v does not lxc type", vmid)
	}

	ErrNotImplements = func(vmid int) error {
		return fmt.Errorf("virtual device with VMID: %v does not implements this method", vmid)
	}

	ErrDidNotImplemented = func(vmid, statusCode int) error {
		return fmt.Errorf("VMID: %v command did not implemented. with status code: %v", vmid, statusCode)
	}

	ErrBadStatusCode = func(code int) error {
		return fmt.Errorf("bad response status code: %v", code)
	}
)
