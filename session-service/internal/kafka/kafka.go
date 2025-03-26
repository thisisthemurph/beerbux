package kafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
)

func EnsureKafkaTopics(brokers []string) error {
	if err := kafkatopic.EnsureTopicExists(brokers, kafka.TopicConfig{
		Topic:             publisher.TopicSessionMemberAdded,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		return fmt.Errorf("failed to ensure session.member.added Kafka topic: %w", err)
	}

	return nil
}
