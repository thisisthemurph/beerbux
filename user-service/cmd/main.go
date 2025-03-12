package main

import (
	"database/sql"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
	"log/slog"
	"net"
	"os"

	"github.com/thisisthemurph/beerbux/user-service/internal/config"
	"github.com/thisisthemurph/beerbux/user-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/internal/server"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	logger.Debug("ensuring Kafka topics exist")
	if err := ensureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		logger.Error("failed to ensure Kafka topics", "error", err)
		return fmt.Errorf("failed to ensure Kafka topics: %w", err)
	}

	logger.Debug("connecting to the database")
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return err
	}
	defer db.Close()

	logger.Debug("connecting to Kafka brokers", "brokers", cfg.Kafka.Brokers)
	userCreatedPublisher := publisher.NewUserCreatedKafkaPublisher(cfg.Kafka.Brokers)
	userUpdatedPublisher := publisher.NewUserUpdatedKafkaPublisher(cfg.Kafka.Brokers)

	queries := user.New(db)
	userServer := server.NewUserServer(queries, userCreatedPublisher, userUpdatedPublisher, logger)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServer(grpcServer, userServer)

	if cfg.Environment.IsDevelopment() {
		reflection.Register(grpcServer)
	}

	logger.Debug("starting user server", "address", cfg.UserServerAddress)
	listener, err := net.Listen("tcp", cfg.UserServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen at %v: %w", cfg.UserServerAddress, err)
	}
	defer listener.Close()

	return grpcServer.Serve(listener)
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
