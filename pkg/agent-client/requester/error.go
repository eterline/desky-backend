package requester

import (
	"errors"
	"fmt"
)

type ClientErr struct {
	err error
}

func (e *ClientErr) Error() string {
	return e.err.Error()
}

var (
	ErrNilConnection = &ClientErr{err: errors.New("connection pointer is nil")}

	ErrResponseNotImplemented = &ClientErr{err: errors.New("response not implemented")}

	ErrBadStatusCode = func(code int) error {
		return &ClientErr{err: fmt.Errorf("bad response status code: %v", code)}
	}
)
