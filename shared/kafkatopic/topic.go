package kafkatopic

import (
	"fmt"
	"github.com/segmentio/kafka-go"
)

func EnsureTopicExists(brokers []string, tc kafka.TopicConfig) error {
	for _, broker := range brokers {
		if err := ensureTopicExistsForBroker(broker, tc); err != nil {
			return err
		}
	}
	return nil
}

func ensureTopicExistsForBroker(broker string, tc kafka.TopicConfig) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("failed to dial broker %v: %w", broker, err)
	}
	defer conn.Close()

	exists, err := topicExists(conn, tc.Topic)
	if err != nil {
		return fmt.Errorf("failed to check if topic %v exists: %w", tc.Topic, err)
	}

	if exists {
		return nil
	}

	if err := conn.CreateTopics(tc); err != nil {
		return fmt.Errorf("failed to create topic %v: %w", tc, err)
	}

	return nil
}

func topicExists(conn *kafka.Conn, topic string) (bool, error) {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return false, err
	}

	for _, p := range partitions {
		if p.Topic == topic {
			return true, nil
		}
	}

	return false, nil
}
