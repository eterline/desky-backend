package handlers

import (
	"fmt"
)

var (
	ErrInterpretationToNumber = func(param string) error {
		return fmt.Errorf("parameter: '%s' can't be interpreted as a number", param)
	}

	ErrEmptyParameter = func(param string) error {
		return fmt.Errorf("parameter: '%s' can't be empty", param)
	}
)
