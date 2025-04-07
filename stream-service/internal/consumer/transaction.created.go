package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/stream-service/internal/sse"
)

type SessionTransactionCreatedKafkaConsumer struct {
	server *sse.Server
}

func NewSessionTransactionCreatedKafkaConsumer(server *sse.Server) KafkaMessageHandler {
	return &SessionTransactionCreatedKafkaConsumer{
		server: server,
	}
}

type SessionTransactionCreatedEvent struct {
	SessionID     string  `json:"session_id"`
	TransactionID string  `json:"transaction_id"`
	CreatorID     string  `json:"creator_id"`
	Total         float64 `json:"total"`
}

func (h SessionTransactionCreatedKafkaConsumer) Handle(ctx context.Context, msg kafka.Message) error {
	var ev SessionTransactionCreatedEvent
	err := json.Unmarshal(msg.Value, &ev)
	if err != nil {
		return err
	}

	data, err := json.Marshal(map[string]string{
		"transactionId": ev.TransactionID,
		"sessionId":     ev.SessionID,
		"creatorId":     ev.CreatorID,
		"total":         fmt.Sprintf("%f", ev.Total),
	})

	message := sse.NewMessage("session.transaction.created", data)
	h.server.BroadcastMessageToRoom(ev.SessionID, message)

	return nil
}
