package controllers

import "errors"

var (
	ErrUncorrectCredentials = errors.New("uncorrect login credentials")
)

var (
	ErrWSNotOpened        = errors.New("websocket did not opened")
	ErrUnknownUnitCommand = errors.New("unknown unit command")
)
