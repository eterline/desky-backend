package sys

import "errors"

var (
	ErrWSNotOpened        = errors.New("websocket did not opened")
	ErrUnknownUnitCommand = errors.New("unknown unit command")
)
