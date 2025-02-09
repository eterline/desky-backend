package handler

import (
	"fmt"
)

type WSHandlerError struct {
	Message  string `json:"message"`
	ClosedWS bool   `json:"close"`
}

func (e *WSHandlerError) Error() string {
	return fmt.Sprintf("(close %v) error: %s", e.ClosedWS, e.Message)
}

func NewWSHandlerError(message string, close bool) *WSHandlerError {
	return &WSHandlerError{
		Message:  message,
		ClosedWS: close,
	}
}

var (
	ErrInterpretationToNumber = func(param string) error {
		return fmt.Errorf("parameter: '%s' can't be interpreted as a number", param)
	}

	ErrEmptyParameter = func(param string) error {
		return fmt.Errorf("parameter: '%s' can't be empty", param)
	}
)
