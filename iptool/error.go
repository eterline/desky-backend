package iptool

import (
	"errors"
	"fmt"
)

type IPToolError struct {
	err error
	ver IPv
}

func (i *IPToolError) Error() string {
	return fmt.Sprintf("iptool error: %s", i.err.Error())
}

var (
	ErrUncompVersion = &IPToolError{
		err: errors.New("uncomplitable ip version"),
	}
)
