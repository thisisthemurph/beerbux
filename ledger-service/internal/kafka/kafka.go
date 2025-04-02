package kafka

import (
	"fmt"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/publisher"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
)

func EnsureKafkaTopics(brokers []string) error {
	topics := []string{
		publisher.TopicLedgerUpdated,
		publisher.TopicLedgerTransactionUpdated,
		publisher.TopicLedgerUserTotalsCalculated,
	}

	for _, topic := range topics {
		if err := kafkatopic.EnsureTopicExists(brokers, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}); err != nil {
			return fmt.Errorf("failed to ensure session.member.added Kafka topic: %w", err)
		}
	}

	return nil
}
