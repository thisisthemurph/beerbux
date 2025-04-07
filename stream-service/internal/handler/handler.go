package handler

import (
	"errors"
	"fmt"
	"github.com/thisisthemurph/beerbux/stream-service/internal/sse"
	"net/http"
)

func getFromQuery(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		http.Error(w, key+" is required", http.StatusBadRequest)
		return "", false
	}
	return value, true
}

func getUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	return getFromQuery(w, r, "user_id")
}

func getSessionID(w http.ResponseWriter, r *http.Request) (string, bool) {
	return getFromQuery(w, r, "session_id")
}

func setServerSentEventHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

type EventStreamWriter struct {
	flusher http.Flusher
	writer  http.ResponseWriter
}

func NewEventStreamWriter(w http.ResponseWriter) (*EventStreamWriter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming unsupported")
	}

	return &EventStreamWriter{
		flusher: flusher,
		writer:  w,
	}, nil
}

func (w *EventStreamWriter) Write(m *sse.Message) error {
	_, err := fmt.Fprint(w.writer, m.String())
	if err != nil {
		return err
	}
	w.flusher.Flush()
	return nil
}
