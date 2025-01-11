package appsfile

import (
	"errors"
	"fmt"
)

type AppsFileServiceErrors struct {
	err error
}

func (a *AppsFileServiceErrors) Error() string {
	return fmt.Sprintf("apps file service error: %s", a.err.Error())
}

func IsAppsFileServiceError(e error) bool {

	if e == nil {
		return false
	}

	_, ok := e.(*AppsFileServiceErrors)
	return ok
}

var (
	ErrCannotOpen = func(path string) *AppsFileServiceErrors {
		return &AppsFileServiceErrors{fmt.Errorf("can't open file '%s'. opening default", path)}
	}

	ErrQueryOutOfRange = errors.New("app query number out of range")
)
