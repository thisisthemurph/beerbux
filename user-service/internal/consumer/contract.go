package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
)

// KafkaReader is an interface that defines the methods required to read messages from Kafka.
// This is used for testing purposes to allow mocking of the Kafka reader.
type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}
