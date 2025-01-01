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

	ErrBadStatusCode = func(code int) error {
		return &NodesErr{fmt.Sprintf("bad response status code: %v", code)}
	}
)
