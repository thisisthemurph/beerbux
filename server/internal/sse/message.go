package sse

import "fmt"

type Message struct {
	Topic string
	Key   string
	Value []byte
}

func NewMessage(topic, key string, value []byte) *Message {
	return &Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}
}

func NewHeartbeatMessage() *Message {
	return &Message{
		Topic: "heartbeat",
		Value: []byte(`{"message": "keep-alive"}`),
	}
}

// String returns the message in the required ServerSentEvent format.
func (m *Message) String() string {
	return fmt.Sprintf("event: %s\ndata: %s\n\n", m.Topic, string(m.Value))
}
