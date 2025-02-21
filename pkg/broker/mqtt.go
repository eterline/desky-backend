package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SenderMQTT struct {
	mq   mqtt.Client
	opts *ClientOptions

	workers []WorkerEntry
	waiter  sync.WaitGroup

	closer context.CancelFunc
	ctx    context.Context
}

func NewSender() *SenderMQTT {
	return NewSenderWithContext(context.Background())
}

func NewSenderWithContext(ctx context.Context, opts ...OptionFunc) *SenderMQTT {

	o := &ClientOptions{mqtt.NewClientOptions(), 0}

	for _, optionSet := range opts {
		optionSet(o)
	}

	context, close := context.WithCancel(ctx)

	sender := &SenderMQTT{
		mq:  mqtt.NewClient(o.ClientOptions),
		ctx: context,

		workers: make([]WorkerEntry, 0),

		closer: close,
		opts:   o,
	}

	return sender
}

func (s *SenderMQTT) Connect(timeout time.Duration) error {

	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	defer close(done)

	go func() {
		token := s.mq.Connect()
		token.Wait()
		done <- token.Error()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("mqtt connect error: %v", err)
		}
		return nil

	case <-ctx.Done():
		return fmt.Errorf("mqtt connection timeout")
	}
}

func (s *SenderMQTT) onConnect() {
	for _, w := range s.workers {
		worker := w
		s.waiter.Add(1)

		fmt.Println("Starting worker...")

		go func() {
			defer s.waiter.Done()
			worker.worker(s.ctx, worker.topic)
			fmt.Println("Worker stopped")
		}()
	}
}

func (s *SenderMQTT) TestEstablish() error {
	if !s.mq.IsConnected() {
		return fmt.Errorf("connection is down")
	}
	return nil
}

func (s *SenderMQTT) Exit() {
	if s.ctx.Err() == nil {
		s.closer()
	}
	s.mq.Disconnect(100)

	s.waiter.Wait()
}

type SenderMQTTTopic struct {
	s *SenderMQTT

	topic    string
	qos      byte
	retained bool
}

func (s *SenderMQTT) InitTopic(topic string) *SenderMQTTTopic {
	mqttTopic := &SenderMQTTTopic{
		s:        s,
		topic:    topic,
		qos:      s.opts.QoS,
		retained: true,
	}

	return mqttTopic
}

func (t *SenderMQTTTopic) ExtendQoS(q QoSValue) {
	t.qos = byte(q)
}

func (t *SenderMQTTTopic) UnRetain() {
	t.retained = false
}

func (t *SenderMQTTTopic) Push(data interface{}) error {
	if err := t.s.TestEstablish(); err != nil {
		return fmt.Errorf("mqtt broker push message error: %v", err)
	}

	tx := t.s.mq.Publish(t.topic, t.qos, t.retained, data)

	select {
	case <-tx.Done():
		if err := tx.Error(); err != nil {
			return fmt.Errorf("mqtt broker push message error: %v", err)
		}
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("mqtt broker push message timeout")
	}
}

func (t *SenderMQTTTopic) PushJSON(v any) error {
	if err := t.s.TestEstablish(); err != nil {
		return fmt.Errorf("mqtt broker push message error: %v", err)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("mqtt broker marshal json message error: %v", err)
	}

	return t.Push(data)
}
