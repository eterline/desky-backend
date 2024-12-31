package ve

import (
	"context"
	"math"

	"github.com/eterline/desky-backend/internal/configuration"
	proxmox "github.com/eterline/desky-backend/pkg/proxm-ve-tool/client"
	nodes "github.com/eterline/desky-backend/pkg/proxm-ve-tool/nodes"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/virtual"
)

const (
	DataDivisorNumber = 1024 // Divisor value for data measures
)

// Implements Proxmox VE sessions provider for web-server
type ProxmoxProvider struct {
	ErrCount     int
	SessionStack []*nodes.NodeProvider
}

func NewProvide() (*ProxmoxProvider, error) {
	config := configuration.GetConfig().Proxmox
	Provider := &ProxmoxProvider{}

	for _, confInstance := range config {
		cfg := proxmox.InitSession(
			confInstance.Username,
			confInstance.Password,
			confInstance.ApiURL,
			"",
			confInstance.SSLCheck,
		)

		ss, err := proxmox.Connect(cfg)
		if err == nil {
			Provider.SessionStack = append(
				Provider.SessionStack,
				nodes.NewNodeProvider(ss),
			)
			continue
		}
		Provider.ErrCount++
	}

	if len(Provider.SessionStack) == 0 {
		return nil, ErrNoValidSessions
	}

	return Provider, nil
}

func (pp *ProxmoxProvider) AnyValidConns() bool {
	return len(pp.SessionStack) > 0
}

func (pp *ProxmoxProvider) AvailSessions() int {
	return len(pp.SessionStack)
}

type ProvideInstance struct {
	*nodes.NodeProvider
}

// returns session by ID from Session Stack
func (pp *ProxmoxProvider) GetSession(sessionID int) (instance *ProvideInstance, err error) {

	if !pp.AnyValidConns() {
		return nil, ErrNoValidSessions
	}

	if (len(pp.SessionStack) - 1) < sessionID {
		return nil, ErrNoSessionWithId
	}

	return &ProvideInstance{
		NodeProvider: pp.SessionStack[sessionID],
	}, nil
}

func (pi *ProvideInstance) ResolveNode(node string) (v *nodes.ProxmoxNode, err error) {
	return pi.NodeInstance(node)
}

func (pi *ProvideInstance) ResolveDevice(node string, vmid int) (v *virtual.VirtMachine, err error) {

	nodeInstance, err := pi.ResolveNode(node)
	if err != nil {
		return nil, err
	}

	return nodeInstance.VirtMachineInstance(vmid)
}

func (pi *ProvideInstance) NodeStatus(ctx context.Context, node string) (status *PVENodeStatus, err error) {

	nodeInstance, err := pi.ResolveNode(node)
	if err != nil {
		return nil, err
	}

	nodeStatus, err := nodeInstance.Status(ctx)
	if err != nil {
		return nil, err
	}

	data := nodeStatus.Data

	return &PVENodeStatus{
		Name:    node,
		AVGLoad: AVGLoadData(data.Loadavg),

		FS: FSData{
			Used:  data.Rootfs.Used / DataDivisorNumber / DataDivisorNumber,
			Total: data.Rootfs.Total / DataDivisorNumber / DataDivisorNumber,
		},

		RAM: RAMData{
			Used:  data.Memory.Used / DataDivisorNumber / DataDivisorNumber,
			Total: data.Memory.Total / DataDivisorNumber / DataDivisorNumber,
		},

		CPU: CPUData{
			Load:      math.Ceil(data.CPU),
			Model:     data.Cpuinfo.Model,
			Cores:     data.Cpuinfo.Cores,
			Frequency: data.Cpuinfo.Mhz,
		},

		Uptime: int64(data.Uptime),
		Kernel: data.CurrentKernel.Release,
	}, nil
}

func (pi *ProvideInstance) DeviceList(ctx context.Context, node string) (listDev *DevicesList, err error) {

	nodeInstance, err := pi.ResolveNode(node)
	if err != nil {
		return nil, err
	}

	lxcList, err := nodeInstance.LXCList(ctx)
	if err != nil {
		return nil, err
	}

	qemuList, err := nodeInstance.VMList(ctx)
	if err != nil {
		return nil, err
	}

	listDev = &DevicesList{}

	for _, prop := range lxcList.Data {
		listDev.LXC = append(
			listDev.LXC,
			TypeDevice{
				Status: prop.Status,
				Name:   prop.Name,
				Tags:   prop.Tags,
				VMID:   prop.Vmid,
				CPUS:   prop.Cpus,
				NetIn:  int64(prop.Netin),
				NetOut: int64(prop.Netout),
				Uptime: int64(prop.Uptime),
				PID:    prop.Pid,
			},
		)
	}

	for _, prop := range qemuList.Data {
		listDev.QEMU = append(
			listDev.QEMU,
			TypeDevice{
				Status: prop.Status,
				Name:   prop.Name,
				Tags:   prop.Tags,
				VMID:   prop.Vmid,
				CPUS:   prop.Cpus,
				NetIn:  int64(prop.Netin),
				NetOut: int64(prop.Netout),
				Uptime: int64(prop.Uptime),
				PID:    prop.Pid,
			},
		)
	}

	return listDev, nil
}
