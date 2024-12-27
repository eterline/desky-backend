package nodes

import (
	"context"

	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/client"
	"github.com/eterline/desky-backend/pkg/proxm-ve-tool/utils"
)

type NodeProvider struct {
	session *client.Session
	nodes   *NodeList
}

func NewNodeProvider(session *client.Session) *NodeProvider {
	return &NodeProvider{
		session: session,
		nodes:   nil,
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

	if 200 > lst.Code || lst.Code > 299 {
		return nil, ErrBadStatusCode(lst.Code)
	}

	err = req.Resolve(&lst)
	if err == nil {
		np.nodes = lst
	}
	return lst, err
}

func (np *NodeProvider) Node(name string) (node *ProxmoxNode, err error) {
	list := &NodeList{}

	if np.nodes == nil {
		if list, err = np.GetNodes(context.Background()); err != nil {
			return nil, err
		}
	}

	if !utils.ContainsInListOfStruct(list.Data, name) {
		return nil, ErrNodeNotExists(name)
	}

	return &ProxmoxNode{
		session: np.session,
		Name:    name,
	}, nil
}
