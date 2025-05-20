package handler

import (
	"beerbux/internal/sse"
	"errors"
	"fmt"
	"net/http"
)

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
