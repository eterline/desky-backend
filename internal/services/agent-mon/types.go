package agentmon

import (
	"context"
	"sync"

	"github.com/eterline/desky-backend/pkg/broker"
)

type CacheStorage interface {
	GetValue(key any) any
	PushValue(key any, value any)
	CleanValue(key any)
}

type BrokerListener interface {
	ListenTopic(topic string, msgHandle func(broker.Message)) error
}

type Provider interface {
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
	done     bool
	mu       sync.Mutex
}

type AgentRequest interface {
	ValueAPI() string
	ValueToken() string
}

type ValidateData struct {
	URL      string
	ID       string
	Hostname string
	Err      error
}
