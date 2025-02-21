package agentmon

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/pkg/broker"
)

type AgentDataMessage struct {
	ID        string                  `json:"host-id"`
	Data      models.AgentStatsObject `json:"data"`
	Timestamp int64                   `json:"timestamp"`
}

type AgentMonitorServiceWithBroker struct {
	broker BrokerListener

	agentStats map[string]AgentDataMessage
	agentStack map[string]models.SessionCredentials

	ctx context.Context
	mu  sync.RWMutex
}

func NewAgentMonitorServiceWithBroker(
	ctx context.Context,
	broker BrokerListener,
) *AgentMonitorServiceWithBroker {
	return &AgentMonitorServiceWithBroker{
		broker: broker,
		ctx:    ctx,

		agentStats: make(map[string]AgentDataMessage),
		agentStack: make(map[string]models.SessionCredentials),
	}
}

func (ab *AgentMonitorServiceWithBroker) RunDataUpdater(topicListen string) error {
	return ab.broker.ListenTopic(topicListen, func(m broker.Message) {
		defer m.Ack()

		data := new(AgentDataMessage)

		if err := json.Unmarshal(m.Payload(), data); err != nil {
			return
		}

		ab.mu.Lock()
		defer ab.mu.Unlock()

		if _, ok := ab.agentStack[data.ID]; !ok {
			ab.agentStack[data.ID] = models.SessionCredentials{
				Hostname: data.Data.Host.Hostname,
				ID:       data.ID,
				URL:      "-",
			}
		}

		data.Timestamp = time.Now().Unix()
		ab.agentStats[data.ID] = *data
	})
}

func (ab *AgentMonitorServiceWithBroker) List() []models.SessionCredentials {

	list := make([]models.SessionCredentials, 0)

	ab.mu.RLock()
	defer ab.mu.RUnlock()

	for _, data := range ab.agentStack {
		list = append(list, data)
	}

	return list
}

func (ab *AgentMonitorServiceWithBroker) Pool() (<-chan any, context.CancelFunc) {

	ch := make(chan any)
	ctx, cancel := context.WithCancel(ab.ctx)

	go ab.collectStatsTo(ctx, ch) // first send after start

	go func() {
		tick := time.NewTicker(5 * time.Second)
		defer tick.Stop()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return

			case <-tick.C:
				ab.collectStatsTo(ctx, ch)
			}
		}
	}()

	return ch, cancel
}

func (ab *AgentMonitorServiceWithBroker) collectStatsTo(ctx context.Context, channel chan any) {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	for key, data := range ab.agentStats {
		if ctx.Err() == nil {
			channel <- data

			if time.Unix(
				data.Timestamp, 0,
			).Add(
				10 * time.Second,
			).After(
				time.Unix(data.Timestamp, 0),
			) {
				delete(ab.agentStats, key)
			}
		}
	}
}
