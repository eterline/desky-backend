package agentmon

import (
	"context"
	"time"

	"github.com/eterline/desky-backend/internal/models"
)

type AgentDataMessage struct {
	ID       string
	Hostname string
	Data     any
}

type AgentWithBroker struct {
	broker       BrokerListener
	ctx          context.Context
	agentStorage map[string]AgentDataMessage
}

func NewAgentWithBroker(
	ctx context.Context,
	broker BrokerListener,
) *AgentWithBroker {
	return &AgentWithBroker{
		broker:       broker,
		ctx:          ctx,
		agentStorage: make(map[string]AgentDataMessage),
	}
}

func (ab *AgentWithBroker) RunDataUpdater() error {
	return nil
}

func (ab *AgentWithBroker) List() []models.SessionCredentials {
	list := make([]models.SessionCredentials, 0)

	for _, data := range ab.agentStorage {
		list = append(list, models.SessionCredentials{
			Hostname: data.Hostname,
			ID:       data.ID,
			URL:      "-",
		})
	}

	return list
}

func (ab *AgentWithBroker) Pool() (<-chan any, context.CancelFunc) {

	ch := make(chan any)
	ctx, cancel := context.WithCancel(ab.ctx)

	go func() {
		tick := time.NewTicker(5 * time.Second)
		defer tick.Stop()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return

			case <-tick.C:
				for _, data := range ab.agentStorage {
					ch <- data
				}
			}
		}
	}()

	return ch, cancel
}
