package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/config"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/kafka"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	if err := run(cfg, logger); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("starting ledger-service", "env", cfg.Environment)

	if err := kafka.EnsureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return fmt.Errorf("failed to ensure Kafka topics: %w", err)
	}

	db, err := connectToDatabase(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	ctx, cancel := setupSignalHandler()
	defer cancel()

	// ledgerUpdatedChan is appended to when a set of ledger items are added
	// to the database from a transaction.created event.
	ledgerUpdatedChan := make(chan []event.LedgerUpdateEvent)

	setupAndRunKafkaConsumers(ctx, cfg, logger, db, ledgerUpdatedChan)
	ledgerUpdatedPublisher := publisher.NewLedgerUpdatedKafkaPublisher(cfg.Kafka.Brokers)
	ledgerTransactionUpdatedPublisher := publisher.NewLedgerTransactionUpdatedKafkaPublisher(cfg.Kafka.Brokers)

	for {
		select {
		case <-ctx.Done():
			logger.Debug("shutting down")
			return nil
		case updates := <-ledgerUpdatedChan:
			logger.Debug("ledger updated", "result", updates)

			if err := ledgerTransactionUpdatedPublisher.Publish(ctx, updates); err != nil {
				logger.Error("failed to publish ledger transaction updated event", "error", err)
			}

			for _, update := range updates {
				if err := ledgerUpdatedPublisher.Publish(ctx, update); err != nil {
					logger.Error("failed to publish ledger updated event", "error", err)
				}
			}
		}
	}
}

func connectToDatabase(driver, uri string) (*sql.DB, error) {
	db, err := sql.Open(driver, uri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrateDatabase(db, driver); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func migrateDatabase(db *sql.DB, driver string) error {
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Up(db, "./internal/db/migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func setupSignalHandler() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	return ctx, cancel
}

func setupAndRunKafkaConsumers(
	ctx context.Context,
	cfg *config.Config,
	logger *slog.Logger,
	db *sql.DB,
	ledgerUpdatedChan chan<- []event.LedgerUpdateEvent,
) {
	ledgerRepository := repository.NewLedgerQueries(db)
	updateLedgerHandler := handler.NewUpdateLedgerHandler(ledgerRepository, ledgerUpdatedChan, logger)
	transactionCreatedConsumer := kafka.NewConsumer(logger, cfg.Kafka.Brokers, "transaction.created", "ledger-service")

	go transactionCreatedConsumer.StartListening(ctx, updateLedgerHandler.Handle)
}
