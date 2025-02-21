package broker

import (
	"context"
)

type WorkerMQTTFunc func(context.Context, *SenderMQTTTopic)

type WorkerEntry struct {
	topic  *SenderMQTTTopic
	worker WorkerMQTTFunc
}

func (s *SenderMQTT) AddTopicWorker(topic *SenderMQTTTopic, work WorkerMQTTFunc) {
	s.workers = append(s.workers, WorkerEntry{topic: topic, worker: work})
}
