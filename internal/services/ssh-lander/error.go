package sshlander

import (
	"fmt"

	"github.com/google/uuid"
)

func NewError(uuid uuid.UUID, msg string) *SSHLanderError {
	return &SSHLanderError{
		Message: msg,
		UUID:    uuid,
	}
}

type SSHLanderError struct {
	Message string
	UUID    uuid.UUID
}

func (e *SSHLanderError) Error() string {
	return fmt.Sprintf("sshlander uuid: %s error: %s", e.UUID.String(), e.Message)
}
