package nodes

import (
	"context"
)

const (
	UpdatePath = "/apt/update"
	VersPath   = "/apt/versions"
)

func (pn *ProxmoxNode) GetAptVersions(ctx context.Context) (versions *AptVersions, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(VersPath))
	defer req.EndTask()
	versions = &AptVersions{}

	versions.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if err := req.Resolve(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func (pn *ProxmoxNode) GetAptUpdates(ctx context.Context) (versions *AptUpdates, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(UpdatePath))
	defer req.EndTask()
	versions = &AptUpdates{}

	versions.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if err := req.Resolve(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func (pn *ProxmoxNode) AptUpgrade(ctx context.Context) (versions *AptUpgrade, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(UpdatePath))
	defer req.EndTask()
	versions = &AptUpgrade{}

	versions.Code, err = req.POST()
	if err != nil {
		return nil, err
	}

	if err := req.Resolve(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}
