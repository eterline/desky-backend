package agentmon

import "context"

type Provider interface {
	IsValid() bool
	Parameter(string) (any, error)
}

type Session struct {
	Hostname string
	ID       string
	URL      string
	Provider
}

type AgentMonitorService struct {
	sessions []Session
	ctx      context.Context
}
