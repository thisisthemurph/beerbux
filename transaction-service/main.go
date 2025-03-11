package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
	"github.com/thisisthemurph/beerbux/transaction-service/internal/config"
	"github.com/thisisthemurph/beerbux/transaction-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/transaction-service/internal/server"
	"github.com/thisisthemurph/beerbux/transaction-service/protos/transactionpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	if err := run(cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	logger.Debug("ensuring Kafka topics exist")
	if err := ensureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return fmt.Errorf("failed to ensure Kafka topics: %w", err)
	}

	sessionClientConn, err := grpc.NewClient(cfg.SessionServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error connecting to session-service: %w", err)
	}

	sessionClient := sessionpb.NewSessionClient(sessionClientConn)
	transactionCreatedPublisher := publisher.NewTransactionCreatedKafkaPublisher(cfg.Kafka.Brokers)
	transactionServer := server.NewTransactionServer(sessionClient, transactionCreatedPublisher)

	gs := grpc.NewServer()
	transactionpb.RegisterTransactionServer(gs, transactionServer)

	logger.Debug("listening for transactions", "address", cfg.TransactionServerAddress)
	listener, err := net.Listen("tcp", cfg.TransactionServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if cfg.Environment.IsDevelopment() {
		reflection.Register(gs)
	}

	if err := gs.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func ensureKafkaTopics(brokers []string) error {
	if err := kafkatopic.EnsureTopicExists(brokers, kafka.TopicConfig{
		Topic:             publisher.TopicTransactionCreated,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		return fmt.Errorf("failed to ensure %v Kafka topic: %w", publisher.TopicTransactionCreated, err)
	}

	return nil
}
