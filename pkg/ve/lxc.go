package ve

import (
	"context"
	"sort"

	"github.com/luthermonson/go-proxmox"
)

type LXC struct {
	MemFmt, DiskFmt, SwapFmt, UptimeFmt string
	Device                              *proxmox.Container
}

func (node *VENode) LXCList() ([]LXC, error) {
	cts, err := node.Containers(node.Context)
	if err != nil {
		return nil, err
	}
	var l []LXC
	for _, i := range cts {
		l = append(l, CollectCT(i))
	}

	sort.Slice(l, func(i, j int) (less bool) {
		return l[i].Device.VMID < l[j].Device.VMID
	})
	return l, nil
}

func (node *VENode) LXCget(id int) (LXC, error) {
	ct, err := node.Container(node.Context, id)
	if err != nil {
		return LXC{}, err
	}
	return CollectCT(ct), nil
}

func (ct LXC) Shutdown() {
	ct.Device.Shutdown(context.Background(), false, 0)
}

func CollectCT(lxc *proxmox.Container) LXC {
	return LXC{
		MemFmt:    sizeStrMB(lxc.MaxMem),
		DiskFmt:   sizeStrMB(lxc.MaxDisk),
		SwapFmt:   sizeStrMB(lxc.MaxSwap),
		UptimeFmt: uptimeStr(lxc.Uptime),
		Device:    lxc,
	}
}
