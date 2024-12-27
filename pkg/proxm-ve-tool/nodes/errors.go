package nodes

import "fmt"

var (
	ErrNodeNotExists = func(name string) error {
		return fmt.Errorf("node '%s' does not exists", name)
	}

	ErrBadStatusCode = func(code int) error {
		return fmt.Errorf("bad response status code: %v", code)
	}
)
