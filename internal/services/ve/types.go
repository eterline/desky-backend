package ve

import nodes "github.com/eterline/desky-backend/pkg/proxm-ve-tool/nodes"

// Implements Proxmox VE sessions provider for web-server
type VEService struct {
	ErrCount     int
	SessionStack []*nodes.NodeProvider
}

type ProvideInstance struct {
	*nodes.NodeProvider
}

type SessionNodes []nodes.NodeUnit
