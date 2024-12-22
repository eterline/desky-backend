package configuration

import "errors"

var (
	ErrUnknownExt = errors.New("unknown configuration file extension type. use default.")
)
