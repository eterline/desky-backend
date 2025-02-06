package authorization

import (
	"errors"
	"fmt"
)

type AuthorizationServiceError struct {
	err error
}

func (ausr *AuthorizationServiceError) Error() string {
	return fmt.Sprintf("authorization service error: %s", ausr.err.Error())
}

var (
	ErrVerifyPassword = &AuthorizationServiceError{
		err: errors.New("verification password error"),
	}
)
