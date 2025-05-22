package handler

import (
	"beerbux/internal/sse"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"log/slog"
	"net/http"
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
	userID, ok := url.Query.GetString(r, "user_id")
	if !ok {
		send.BadRequest(w, "user_id is required")
		return
	}
	sessionID, ok := url.Query.GetString(r, "session_id")
	if !ok {
		send.BadRequest(w, "session_id is required")
		return
	}

	setServerSentEventHeaders(w)

	eventStreamWriter, err := NewEventStreamWriter(w)
	if err != nil {
		if errors.Is(err, ErrStreamingUnsupported) {
			send.InternalServerError(w, "Streaming unsupported")
		} else {
			h.logger.Error("Error creating stream writer", "error", err)
			send.InternalServerError(w, "Could not connect to streaming")
		}
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
			if err := eventStreamWriter.Write(message); err != nil {
				h.logger.Error("Failed to write SSE message", "error", err)
				room.RemoveClient(userID)
				return
			}
		}
	}
}
