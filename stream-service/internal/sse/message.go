package sse

import "fmt"

type Message struct {
	Event string
	Data  []byte
}

func NewMessage(event string, data []byte) *Message {
	return &Message{
		Event: event,
		Data:  data,
	}
}

func NewHeartbeatMessage() *Message {
	return &Message{
		Event: "heartbeat",
		Data:  []byte(`{"message": "keep-alive"}`),
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("event: %s\ndata: %s\n\n", m.Event, string(m.Data))
}
