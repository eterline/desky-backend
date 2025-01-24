package configuration

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownExt = errors.New("unknown configuration file extension type. use default.")
	ErrRead       = func(err error) error {
		return fmt.Errorf("can't read config file: %s", err.Error())
	}
)
