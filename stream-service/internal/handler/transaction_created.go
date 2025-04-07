package handler

import (
	"log/slog"
	"net/http"

	"github.com/thisisthemurph/beerbux/stream-service/internal/sse"
)

type SessionTransactionCreatedHandler struct {
	Server *sse.Server
	logger *slog.Logger
}

func NewSessionTransactionCreatedHandler(logger *slog.Logger, server *sse.Server) http.Handler {
	return &SessionTransactionCreatedHandler{
		Server: server,
		logger: logger,
	}
}

func (h *SessionTransactionCreatedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(w, r)
	if !ok {
		return
	}

	sessionID, ok := getSessionID(w, r)
	if !ok {
		return
	}

	setServerSentEventHeaders(w)

	eventStreamWriter, err := NewEventStreamWriter(w)
	if err != nil {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	room := h.Server.GetOrCreateRoom(sessionID)
	client := sse.NewClient(userID)
	room.AddClient(client)

	notify := r.Context().Done()
	go func() {
		<-notify
		room.RemoveClient(userID)
	}()

	for {
		select {
		case <-client.Done:
			return
		case message := <-client.Ch:
			h.logger.Info("Sending SSE message", "message", message)
			if err := eventStreamWriter.Write(message); err != nil {
				h.logger.Error("Failed to write SSE message", "error", err)
				room.RemoveClient(userID)
				return
			}
		}
	}
}
