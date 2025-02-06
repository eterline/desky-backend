package agentmon

import (
	"context"
	"time"

	"github.com/eterline/desky-backend/internal/models"
)

func New(ctx context.Context) *AgentMonitorService {
	return &AgentMonitorService{
		sessions: make([]Session, 0),
		ctx:      ctx,
		done:     false,
	}
}

func (a *AgentMonitorService) AddSession(p Provider, hostname, id, url string) {
	a.sessions = append(a.sessions, Session{
		Hostname: hostname,
		ID:       id,
		Provider: p,
		URL:      url,
	})
}

func (a *AgentMonitorService) List() (data []models.SessionCredentials) {
	data = make([]models.SessionCredentials, len(a.sessions))

	for i, s := range a.sessions {
		data[i] = models.SessionCredentials{
			Hostname: s.Hostname,
			ID:       s.ID,
			URL:      s.URL,
		}
	}

	return
}

func (a *AgentMonitorService) Pool() (ch chan models.FetchedResponse, cancel context.CancelFunc) {
	ch = make(chan models.FetchedResponse)
	ctx, cancel := context.WithCancel(a.ctx)

	go func(ch chan<- models.FetchedResponse) {
		tick := time.NewTicker(5 * time.Second)
		defer tick.Stop()

		for {
			select {

			case <-ctx.Done():
				close(ch)
				return

			case <-tick.C:
				for _, session := range a.sessions {
					go func(ch chan<- models.FetchedResponse) {
						ch <- fetchAll(session)
					}(ch)
				}

			}
		}

	}(ch)

	return
}

func fetchAll(s Session) models.FetchedResponse {

	data := models.FetchedResponse{
		ID:   s.ID,
		Data: make(map[string]any),
	}

	for _, export := range models.ExporterList {
		info, err := s.Parameter(export)
		if info == nil || err != nil {
			continue
		}

		data.Data[export] = info
	}

	return data
}

func fetchAllToChannel(s Session, ch chan<- models.FetchedResponse) {
	data := fetchAll(s)
	ch <- data
}
