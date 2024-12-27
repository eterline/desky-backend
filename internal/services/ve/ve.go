package ve

import (
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

func (pp *ProxmoxProvider) ResolveDevice(sessionID int, node string, vmid int) (v *virtual.VirtMachine, err error) {

	if !pp.AnyValidConns() {
		return nil, ErrNoValidSessions
	}

	provider := pp.Providers[sessionID]

	nodeInstance, err := provider.Node(node)
	if err != nil {
		return nil, err
	}

	return nodeInstance.VirtMachineInstance(vmid)
}
