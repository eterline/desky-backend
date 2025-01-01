package virtual

import "fmt"

type VirtualErr struct {
	string
}

func (e *VirtualErr) Error() string {
	return e.string
}

var (
	ErrVirtualNotExists = func(vmid int) error {
		return &VirtualErr{fmt.Sprintf("virtual device with VMID: %v does not exists", vmid)}
	}

	ErrNotQEMU = func(vmid int) error {
		return &VirtualErr{fmt.Sprintf("virtual device with VMID: %v does not qemu type", vmid)}
	}

	ErrNotLXC = func(vmid int) error {
		return &VirtualErr{fmt.Sprintf("virtual device with VMID: %v does not lxc type", vmid)}
	}

	ErrNotImplements = func(vmid int) error {
		return &VirtualErr{fmt.Sprintf("virtual device with VMID: %v does not implements this method", vmid)}
	}

	ErrDidNotImplemented = func(vmid, statusCode int) error {
		return &VirtualErr{fmt.Sprintf("VMID: %v command did not implemented. with status code: %v", vmid, statusCode)}
	}

	ErrBadStatusCode = func(code int) error {
		return &VirtualErr{fmt.Sprintf("bad response status code: %v", code)}
	}
)
