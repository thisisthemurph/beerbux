package handler

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type KafkaMessageHandler interface {
	Handle(ctx context.Context, msg kafka.Message) error
}
