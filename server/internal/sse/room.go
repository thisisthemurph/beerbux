package sse

import "sync"

type Room struct {
	id      string
	clients map[string]*Client
	mu      sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		id:      id,
		clients: make(map[string]*Client),
	}
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client.id] = client
}

func (r *Room) RemoveClient(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if client, ok := r.clients[clientID]; ok {
		delete(r.clients, clientID)
		client.Close()
	}
}
