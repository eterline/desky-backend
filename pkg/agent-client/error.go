package agentclient

import "errors"

type AgentErr struct {
	err error
}

func (e *AgentErr) Error() string {
	return e.err.Error()
}

var (
	ErrForbidden         = &AgentErr{err: errors.New("invalid api key")}
	ErrExporterNotExists = &AgentErr{err: errors.New("invalid exporter type value")}
	ErrInvalidAgent      = &AgentErr{err: errors.New("invalid agent session")}
)
