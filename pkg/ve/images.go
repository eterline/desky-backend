package ve

import (
	"context"
)

func (n *VENode) NodeStorage() string {

	storage, err := n.StorageImages(context.Background())
	if err != nil {
		return ""
	}
	return storage.Content
}
