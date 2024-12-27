package ve

import (
	"context"

	"github.com/eterline/desky-backend/internal/configuration"
	proxmox "github.com/eterline/desky-backend/pkg/proxm-ve-tool/client"
	nodes "github.com/eterline/desky-backend/pkg/proxm-ve-tool/nodes"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/virtual"
)

type ProxmoxProvider struct {
	ErrCount  int
	Providers []*nodes.NodeProvider
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
		if err != nil {
			Provider.ErrCount++
		}

		Provider.Providers = append(
			Provider.Providers,
			nodes.NewNodeProvider(ss),
		)
	}

	if len(Provider.Providers) == 0 {
		return nil, ErrNoValidSessions
	}

	return Provider, nil
}

func (pp *ProxmoxProvider) AnyValidConns() bool {
	return len(pp.Providers) > 0
}

type ProvideInstance struct {
	*nodes.NodeProvider
}

func (pp *ProxmoxProvider) GetSession(sessionID int) (instance *ProvideInstance, err error) {

	if !pp.AnyValidConns() {
		return nil, ErrNoValidSessions
	}

	if (len(pp.Providers) - 1) < sessionID {
		return nil, ErrNoSessionWithId
	}

	return &ProvideInstance{
		NodeProvider: pp.Providers[sessionID],
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
			Used:  data.Rootfs.Used,
			Total: data.Rootfs.Total,
		},
		RAM: RAMData{
			Used:  data.Memory.Used,
			Total: data.Memory.Total,
		},
		CPU: CPUData{
			Load: data.CPU,
		},
	}, nil
}
