package virtual

import (
	"context"
	"fmt"

	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/client"
)

type VirtMachine struct {
	VMID       int
	session    *client.Session
	devType    string
	parentNode string
}

func NewVirt(vmid int, s *client.Session, t string, parent string) *VirtMachine {
	return &VirtMachine{
		VMID:       vmid,
		session:    s,
		devType:    t,
		parentNode: parent,
	}
}

func (v *VirtMachine) urlWithBase(path string) string {
	fmt.Println(fmt.Sprintf("/nodes/%s/%s/%v%s", v.parentNode, v.devType, v.VMID, path))
	return fmt.Sprintf("/nodes/%s/%s/%v%s", v.parentNode, v.devType, v.VMID, path)
}

func (v *VirtMachine) IsQEMU() bool {
	return v.devType == VirtTypeQEMU
}

func (v *VirtMachine) IsLXC() bool {
	return v.devType == VirtTypeLXC
}

func (v *VirtMachine) VirtType() string {
	return v.devType
}

func (v *VirtMachine) Status(ctx context.Context) (status interface{}, err error) {
	req := v.session.MakeRequest(ctx, v.urlWithBase("/status/current"))
	defer req.EndTask()

	code, err := req.GET()
	if err != nil {
		return nil, err
	}

	if 200 > code || code > 299 {
		return nil, ErrBadStatusCode(code)
	}

	switch v.devType {
	case VirtTypeLXC:
		status = &LXCStatus{Code: code}
	case VirtTypeQEMU:
		status = &QEMUStatus{Code: code}
	default:
		panic("unknown virt dev type")
	}

	if err := req.Resolve(status); err != nil {
		return nil, err
	}

	return status, nil
}

func (v *VirtMachine) Start(ctx context.Context) error {
	req := v.session.MakeRequest(ctx, v.urlWithBase("/status/start"))
	defer req.EndTask()

	code, err := req.POST()
	if err != nil {
		return err
	}

	if code != 200 {
		return ErrDidNotImplemented(v.VMID, code)
	}

	return nil
}

func (v *VirtMachine) Shutdown(ctx context.Context) error {
	req := v.session.MakeRequest(ctx, v.urlWithBase("/status/shutdown"))
	defer req.EndTask()

	code, err := req.POST()
	if err != nil {
		return err
	}

	if code != 200 {
		return ErrDidNotImplemented(v.VMID, code)
	}

	return nil
}

func (v *VirtMachine) Stop(ctx context.Context) error {
	req := v.session.MakeRequest(ctx, v.urlWithBase("/status/stop"))
	defer req.EndTask()

	code, err := req.POST()
	if err != nil {
		return err
	}

	if code != 200 {
		return ErrDidNotImplemented(v.VMID, code)
	}

	return nil
}

func (v *VirtMachine) Suspend(ctx context.Context) error {
	if v.IsLXC() {
		return ErrNotImplements(v.VMID)
	}

	req := v.session.MakeRequest(ctx, v.urlWithBase("/status/suspend"))
	defer req.EndTask()

	code, err := req.POST()
	if err != nil {
		return err
	}

	if code != 200 {
		return ErrDidNotImplemented(v.VMID, code)
	}

	return nil
}

func (v *VirtMachine) Resume(ctx context.Context) error {
	if v.IsLXC() {
		return ErrNotImplements(v.VMID)
	}

	req := v.session.MakeRequest(ctx, v.urlWithBase("/status/resume"))
	defer req.EndTask()

	code, err := req.POST()
	if err != nil {
		return err
	}

	if code != 200 {
		return ErrDidNotImplemented(v.VMID, code)
	}

	return nil
}
