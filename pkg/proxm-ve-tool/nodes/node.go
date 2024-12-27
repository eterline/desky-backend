package nodes

import (
	"context"

	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/client"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/utils"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/virtual"
)

const (
	NodesPath = "/nodes/"

	StatusPath  = "/status"
	HostPath    = "/hosts"
	DNSPath     = "/dns"
	APLInfoPath = "/aplinfo"

	QEMUPath = "/qemu"
	LXCPath  = "/lxc"
)

type ProxmoxNode struct {
	session *client.Session
	Name    string
}

func (pn *ProxmoxNode) urlWithName(path string) string {
	return NodesPath + pn.Name + path
}

func (pn *ProxmoxNode) Status(ctx context.Context) (status *NodeStatus, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(StatusPath))
	defer req.EndTask()
	status = &NodeStatus{}

	status.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if 200 > status.Code || status.Code > 299 {
		return nil, ErrBadStatusCode(status.Code)
	}

	if err := req.Resolve(&status); err != nil {
		return nil, err
	}

	return status, nil
}

func (pn *ProxmoxNode) HostsFile(ctx context.Context) (hosts *HostsFile, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(HostPath))
	defer req.EndTask()
	hosts = &HostsFile{}

	hosts.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if 200 > hosts.Code || hosts.Code > 299 {
		return nil, ErrBadStatusCode(hosts.Code)
	}

	if err := req.Resolve(&hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}

func (pn *ProxmoxNode) DNSInfo(ctx context.Context) (dns *DNS, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(DNSPath))
	defer req.EndTask()
	dns = &DNS{}

	dns.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if 200 > dns.Code || dns.Code > 299 {
		return nil, ErrBadStatusCode(dns.Code)
	}

	if err := req.Resolve(&dns); err != nil {
		return nil, err
	}

	return dns, nil
}

func (pn *ProxmoxNode) AplInfo(ctx context.Context) (apl *AplInfo, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(APLInfoPath))
	defer req.EndTask()
	apl = &AplInfo{}

	apl.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if 200 > apl.Code || apl.Code > 299 {
		return nil, ErrBadStatusCode(apl.Code)
	}

	if err := req.Resolve(&apl); err != nil {
		return nil, err
	}

	return apl, nil
}

func (pn *ProxmoxNode) LXCList(ctx context.Context) (lxcs *LXCList, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(LXCPath))
	defer req.EndTask()
	lxcs = &LXCList{}

	lxcs.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if 200 > lxcs.Code || lxcs.Code > 299 {
		return nil, ErrBadStatusCode(lxcs.Code)
	}

	if err := req.Resolve(&lxcs); err != nil {
		return nil, err
	}

	return lxcs, nil
}

func (pn *ProxmoxNode) VMList(ctx context.Context) (vms *VMList, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(QEMUPath))
	defer req.EndTask()
	vms = &VMList{}

	vms.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if 200 > vms.Code || vms.Code > 299 {
		return nil, ErrBadStatusCode(vms.Code)
	}

	if err := req.Resolve(&vms); err != nil {
		return nil, err
	}

	return vms, nil
}

func (pn *ProxmoxNode) VirtMachineInstance(vmid int) (v *virtual.VirtMachine, err error) {

	ctx := context.Background()

	lxcs, err := pn.LXCList(ctx)
	if err != nil {
		return nil, err
	}

	if utils.ContainsInListOfStruct(lxcs.Data, vmid) {
		return virtual.NewVirt(
			vmid, pn.session, virtual.VirtTypeLXC, pn.Name,
		), nil
	}

	vms, err := pn.VMList(ctx)
	if err != nil {
		return nil, err
	}

	if utils.ContainsInListOfStruct(vms.Data, vmid) {
		return virtual.NewVirt(
			vmid, pn.session, virtual.VirtTypeQEMU, pn.Name,
		), nil
	}

	return nil, virtual.ErrVirtualNotExists(vmid)
}
