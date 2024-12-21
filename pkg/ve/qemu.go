package ve

import (
	"context"
	"sort"

	"github.com/luthermonson/go-proxmox"
)

type VM struct {
	MemFmt, DiskFmt, UptimeFmt string
	Device                     *proxmox.VirtualMachine
}

func (node *VENode) VMList() ([]VM, error) {
	vms, err := node.VirtualMachines(node.Context)
	if err != nil {
		return nil, err
	}
	var l []VM
	for _, i := range vms {
		l = append(l, CollectVM(i))
	}

	sort.Slice(l, func(i, j int) (less bool) {
		return l[i].Device.VMID < l[j].Device.VMID
	})
	return l, nil
}

func (node *VENode) VMget(id int) (VM, error) {
	vm, err := node.VirtualMachine(node.Context, id)
	if err != nil {
		return VM{}, err
	}
	return CollectVM(vm), nil
}

func (vm VM) Shutdown() {
	vm.Device.Shutdown(context.Background())
}

func CollectVM(vm *proxmox.VirtualMachine) VM {
	return VM{
		MemFmt:    sizeStrMB(vm.MaxMem),
		DiskFmt:   sizeStrMB(vm.MaxDisk),
		UptimeFmt: uptimeStr(vm.Uptime),
		Device:    vm,
	}
}
