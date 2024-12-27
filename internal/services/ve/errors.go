package ve

import "errors"

var (
	ErrNoValidSessions = errors.New("no one valid session has been connected to servers")
)
