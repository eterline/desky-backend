package nodes

import (
	"context"

	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/client"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/utilities"
)

type NodeProvider struct {
	session *client.Session
}

func NewNodeProvider(session *client.Session) *NodeProvider {
	return &NodeProvider{
		session: session,
	}
}

func (np *NodeProvider) GetNodes(ctx context.Context) (lst *NodeList, err error) {
	lst = &NodeList{}

	req := np.session.MakeRequest(ctx, "/nodes")
	defer req.EndTask()

	lst.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	err = req.Resolve(&lst)
	return lst, err
}

func (np *NodeProvider) NodeInstance(name string) (node *ProxmoxNode, err error) {
	list := &NodeList{}

	if list, err = np.GetNodes(context.Background()); err != nil {
		return nil, err
	}

	if !utilities.ContainsInListOfStruct(list.Data, name) {
		return nil, ErrNodeNotExists(name)
	}

	return &ProxmoxNode{
		session: np.session,
		Name:    name,
	}, nil
}
