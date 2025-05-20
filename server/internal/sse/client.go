package sse

import (
	"fmt"
	"time"
)

var ErrorUnresponsiveClient = fmt.Errorf("client is unresponsive")

type Client struct {
	id   string
	Ch   chan *Message
	Done chan struct{}
}

func NewClient(id string) *Client {
	return &Client{
		id:   id,
		Ch:   make(chan *Message, 16),
		Done: make(chan struct{}),
	}
}

func (c *Client) Send(message *Message) error {
	select {
	case c.Ch <- message:
	case <-time.After(time.Second):
		return ErrorUnresponsiveClient
	}

	return nil
}

func (c *Client) Heartbeat() error {
	select {
	case c.Ch <- NewHeartbeatMessage():
	case <-time.After(time.Second):
		return ErrorUnresponsiveClient
	}

	return nil
}

func (c *Client) Close() {
	select {
	case <-c.Done:
	default:
		close(c.Done)
		close(c.Ch)
	}
}
