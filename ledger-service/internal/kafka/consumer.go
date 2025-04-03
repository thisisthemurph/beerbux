package kafka

import (
	"context"
	"log"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type ConsumerListener interface {
	StartListening(ctx context.Context, handler func(context.Context, kafka.Message) error)
}

type Consumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewConsumer(logger *slog.Logger, brokers []string, topic, groupID string) ConsumerListener {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		logger: logger,
	}
}

func (c *Consumer) StartListening(ctx context.Context, handler func(context.Context, kafka.Message) error) {
	defer c.reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error("failed to read message", "topic", c.reader.Config().Topic, "error", err)
				log.Println("failed to read message", err)
				continue
			}

			c.logger.Debug("Received message", "topic", c.reader.Config().Topic, "offset", msg.Offset)
			if err := handler(ctx, msg); err != nil {
				log.Println("failed to handle message", err)
			}
		}
	}
}
