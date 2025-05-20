package sse

import (
	"errors"
	"log/slog"
	"sync"
)

// Server represents a server that manages multiple rooms and clients.
//
// It is responsible for:
//   - Adding and removing rooms.
//   - Broadcasting messages to all clients in a room.
//   - Sending messages to a specific client in a room.
type Server struct {
	Rooms map[string]*Room
	mu    sync.RWMutex

	logger *slog.Logger
}

func NewServer(logger *slog.Logger) *Server {
	return &Server{
		Rooms:  make(map[string]*Room),
		logger: logger,
	}
}

func (s *Server) GetOrCreateRoom(roomID string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.Rooms[roomID]
	if !ok {
		room = NewRoom(roomID)
		s.Rooms[roomID] = room
	}
	return room
}

func (s *Server) RemoveRoom(roomID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room, ok := s.Rooms[roomID]; ok {
		room.mu.Lock()
		defer room.mu.Unlock()
		for _, client := range room.clients {
			client.Close()
		}
	}
	delete(s.Rooms, roomID)
}

// Heartbeat sends a heartbeat message to all clients in all rooms.
// If a client is unresponsive, it will be removed from the room.
func (s *Server) Heartbeat() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, room := range s.Rooms {
		room.mu.RLock()
		for _, client := range room.clients {
			err := client.Heartbeat()
			if errors.Is(err, ErrorUnresponsiveClient) {
				s.handleUnresponsiveClient(client, room)
			}
		}
		room.mu.RUnlock()
	}
}

func (s *Server) BroadcastMessageToRoom(roomID string, message *Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if room, ok := s.Rooms[roomID]; ok {
		room.mu.RLock()
		defer room.mu.RUnlock()

		for _, client := range room.clients {
			err := client.Send(message)
			if errors.Is(err, ErrorUnresponsiveClient) {
				s.handleUnresponsiveClient(client, room)
			}
		}
	}
}

func (s *Server) SendMessageToClient(roomID, clientID string, message *Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if room, roomOK := s.Rooms[roomID]; roomOK {
		room.mu.RLock()
		defer room.mu.RUnlock()

		if client, ok := room.clients[clientID]; ok {
			err := client.Send(message)
			if errors.Is(err, ErrorUnresponsiveClient) {
				s.handleUnresponsiveClient(client, room)
			}
		}
	}
}

// handleUnresponsiveClient handles the case when a client is unresponsive.
// It removes the client from the room and closes the connection.
// If the room has no clients left, it removes the room from the server.
func (s *Server) handleUnresponsiveClient(client *Client, room *Room) {
	s.logger.Debug("Client unresponsive", "clientID", client.id, "roomID", room.id)

	if len(room.clients) <= 1 {
		s.RemoveRoom(room.id)
	} else {
		room.RemoveClient(client.id)
	}
}
