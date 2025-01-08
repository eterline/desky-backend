package client

import (
	"errors"
	"fmt"
)

type ClientErr struct {
	string
}

func (e *ClientErr) Error() string {
	return e.string
}

var (
	ErrNilConnection = errors.New("connection pointer is nil")

	ErrBadStatusCode = func(code int) error {
		return &ClientErr{fmt.Sprintf("bad response status code: %v", code)}
	}
)
