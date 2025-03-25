package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/config"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()
	if err := run(cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
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

	logger.Debug("connecting to the database")
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		close(doneChan)
	}()

	// ledgerChan is appended to when a new ledger item is added to the database.
	ledgerChan := make(chan handler.MemberTransaction)
	defer close(ledgerChan)

	// ledgerTransactionUpdatedChan is appended to when a transaction has been calculated in the ledger table.
	ledgerTransactionUpdatedChan := make(chan []handler.MemberTransaction)
	defer close(ledgerTransactionUpdatedChan)

	// errChan is appended to when an error occurs while processing a transaction.created event.
	errChan := make(chan event.TransactionCreatedEvent)
	defer close(errChan)

	ledgerRepository := repository.NewLedgerQueries(db)
	ledgerUpdatedPublisher := publisher.NewLedgerUpdatedKafkaPublisher(cfg.Kafka.Brokers)
	ledgerTransactionUpdatedPublisher := publisher.NewLedgerTransactionUpdatedKafkaPublisher(cfg.Kafka.Brokers)
	updateLedgerHandler := handler.NewUpdateLedgerHandler(ledgerRepository, logger)

	logger.Debug("listening for transaction.created events")
	go listenForTransactionCreatedEvents(
		ctx,
		logger,
		cfg.Kafka.Brokers,
		updateLedgerHandler,
		ledgerChan,
		ledgerTransactionUpdatedChan,
		errChan)

	for {
		select {
		case l := <-ledgerChan:
			logger.Debug("ledger updated", "result", l)
			err := ledgerUpdatedPublisher.Publish(ctx, event.LedgerUpdatedEvent{
				ID:            l.ID,
				TransactionID: l.TransactionID,
				SessionID:     l.SessionID,
				UserID:        l.UserID,
				ParticipantID: l.ParticipantID,
				Amount:        l.Amount,
			})
			if err != nil {
				logger.Error("failed to publish ledger updated event", "error", err)
			}
		case lt := <-ledgerTransactionUpdatedChan:
			logger.Debug("ledger transaction updated", "result", lt)
			err := ledgerTransactionUpdatedPublisher.Publish(ctx, lt)
			if err != nil {
				logger.Error("failed to publish ledger transaction updated event", "error", err)
			}
		case errEv := <-errChan:
			logger.Error("failed to process transaction", "transactionID", errEv.TransactionID)
		case <-doneChan:
			logger.Debug("shutting down")
			cancel()
			return nil
		}
	}
}

func ensureKafkaTopics(brokers []string) error {
	topics := []string{publisher.TopicLedgerUpdated, publisher.TopicLedgerTransactionUpdated}

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

func listenForTransactionCreatedEvents(
	ctx context.Context,
	logger *slog.Logger,
	brokers []string,
	updateLedgerHandler *handler.UpdateLedgerHandler,
	ledgerChan chan<- handler.MemberTransaction,
	ledgerTransactionUpdatedChan chan<- []handler.MemberTransaction,
	errChan chan<- event.TransactionCreatedEvent,
) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "transaction.created",
		GroupID: "ledger-service",
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			logger.Debug("Kafka consumer shutting down")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				logger.Error("failed to read message", "error", err)
				continue
			}
			logger.Info("received message", "message", msg)

			var ev event.TransactionCreatedEvent
			if err = json.Unmarshal(msg.Value, &ev); err != nil {
				logger.Error("failed to unmarshal event", "error", err)
				continue
			}

			memberTransactions, err := updateLedgerHandler.Handle(ctx, ev)
			if err != nil {
				errChan <- ev
				continue
			}

			ledgerTransactionUpdatedChan <- memberTransactions
			for _, mt := range memberTransactions {
				ledgerChan <- mt
			}
		}
	}
}
