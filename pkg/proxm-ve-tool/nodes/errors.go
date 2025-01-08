package nodes

import "fmt"

type NodesErr struct {
	string
}

func (e *NodesErr) Error() string {
	return e.string
}

var (
	ErrNodeNotExists = func(name string) error {
		return &NodesErr{fmt.Sprintf("node '%s' does not exists", name)}
	}

	ErrDiskPathNotExists = func(path string) error {
		return &NodesErr{fmt.Sprintf("disk with dev-path '%s' does not exists", path)}
	}
)
