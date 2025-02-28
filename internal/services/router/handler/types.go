package handler

import (
	"fmt"
	"net/http"
)

// ====================== Handler types ======================
type APIHandlerFunc func(http.ResponseWriter, *http.Request) (op string, err error)

type BasicHandlerGroup struct{}

// ====================== Parsing types ======================
type (
	ParamOptions struct {
		Numered  []string
		Stringed []string
	}

	ParamFuncOption func(c *ParamOptions)

	Params struct {
		ints map[string]int
		strs map[string]string
	}
)

// ====================== Response types ======================
type APIResponse struct {
	StatusCode int `json:"code"`
	Message    any `json:"message"`
}

type APIErrorResponse struct {
	APIResponse
}

func (ar *APIErrorResponse) Error() string {
	return fmt.Sprintf("api error code: %d. error: %v", ar.StatusCode, ar.Message)
}

type DataErrors map[string]error

func (ar DataErrors) Error() string {

	validatorErrString := "validation error: "

	for key, val := range ar {
		validatorErrString += fmt.Sprintf("%s:%s\n", key, val.Error())
	}

	return validatorErrString
}

// ====================== WebSocket types ======================

type WebSocketLogger interface {
	Info(args ...interface{})
	Error(args ...interface{})

	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type MessageFormat string

const (
	FormatBase64 MessageFormat = "base64"
	FormatJSON   MessageFormat = "base64"
)

type WsSetupMessage struct {
	UUID   string        `json:"uuid"`
	Format MessageFormat `json:"format"`
}

type SocketMessage struct {
	Type int
	Body []byte
}

func (m *SocketMessage) String() string {
	return string(m.Body)
}
