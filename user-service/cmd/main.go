package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
	"github.com/thisisthemurph/beerbux/user-service/internal/application"
	"github.com/thisisthemurph/beerbux/user-service/internal/config"
	"github.com/thisisthemurph/beerbux/user-service/internal/publisher"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	if err := run(logger, cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, cfg *config.Config) error {
	logger.Debug("ensuring Kafka topics exist")
	if err := ensureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		logger.Error("failed to ensure Kafka topics", "error", err)
		return fmt.Errorf("failed to ensure Kafka topics: %w", err)
	}

	app, err := application.New(logger, cfg)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	defer app.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	serverErrChan := make(chan error, 1)

	go app.RunLedgerUpdatedEventConsumer(ctx)

	go func() {
		serverErrChan <- app.RunGRPCServer()
	}()

	select {
	case <-sigChan:
		logger.Info("shutting down")
		return nil
	case err := <-serverErrChan:
		if err != nil {
			logger.Error("gRPC server error", "error", err)
		}
		return err
	}
}

func ensureKafkaTopics(brokers []string) error {
	topics := []string{publisher.TopicUserCreated, publisher.TopicUserUpdated}

	for _, topic := range topics {
		if err := kafkatopic.EnsureTopicExists(brokers, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}); err != nil {
			return fmt.Errorf("failed to ensure %v Kafka topic: %w", topic, err)
		}
	}

	return nil
}
