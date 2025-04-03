package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
	"github.com/thisisthemurph/beerbux/user-service/internal/publisher"
)

func EnsureKafkaTopics(brokers []string) error {
	topics := []string{publisher.TopicUserCreated, publisher.TopicUserUpdated}

	for _, topic := range topics {
		if err := kafkatopic.EnsureTopicExists(brokers, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}); err != nil {
			return fmt.Errorf("failed to ensure %s Kafka topic: %w", topic, err)
		}
	}

	return nil
}
