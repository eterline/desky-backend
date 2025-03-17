package broker

import (
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ListenerMQTT struct {
	mq     mqtt.Client
	opts   *ClientOptions
	closer context.CancelFunc
	ctx    context.Context
}

func NewListener(opts ...OptionFunc) *ListenerMQTT {
	return NewListenerWithContext(context.Background(), opts...)
}

func NewListenerWithContext(ctx context.Context, opts ...OptionFunc) *ListenerMQTT {

	o := &ClientOptions{mqtt.NewClientOptions(), 0}

	for _, optionSet := range opts {
		optionSet(o)
	}

	context, close := context.WithCancel(ctx)

	return &ListenerMQTT{
		mq:     mqtt.NewClient(o.ClientOptions),
		opts:   o,
		closer: close,
		ctx:    context,
	}
}

type Message interface {
	Duplicate() bool
	Qos() byte
	Retained() bool
	Topic() string
	MessageID() uint16
	Payload() []byte
	Ack()
}

func (l *ListenerMQTT) ListenTopic(topic string, msgHandle func(Message)) error {
	token := l.mq.Subscribe(topic, l.opts.QoS, func(c mqtt.Client, m mqtt.Message) {
		msgHandle(m)
	})

	token.Wait()
	return token.Error()
}

func (s *ListenerMQTT) Connect(timeout time.Duration) error {

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
func (s *ListenerMQTT) Connected() bool {
	return s.mq.IsConnected()
}

func (s *ListenerMQTT) Close() {
	s.mq.Disconnect(0)
}
