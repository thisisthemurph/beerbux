package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/config"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
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

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("error connecting to NATS: %w", err)
	}
	defer nc.Close()

	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		close(done)
	}()

	ledgerChan := make(chan *handler.LedgerUpdateResult)
	defer close(ledgerChan)

	errChan := make(chan event.TransactionCreatedEvent)
	defer close(errChan)

	ledgerRepository := repository.NewLedgerQueries(db)
	ledgerUpdatedPublisher := publisher.NewLedgerUpdatedNatsPublisher(nc)
	updateLedgerHandler := handler.NewUpdateLedgerHandler(ledgerRepository, logger)

	msgHandler := transactionCreatedMsgHandlerFn(logger, updateLedgerHandler, ledgerChan, errChan)
	_, err = nc.Subscribe("transaction.created", msgHandler)
	if err != nil {
		return fmt.Errorf("error subscribing to transaction.created: %w", err)
	}

	for {
		select {
		case l := <-ledgerChan:
			logger.Debug("ledger updated", "result", l)
			if err := ledgerUpdatedPublisher.Publish(l.ID, l.TransactionID, l.SessionID, l.UserID, l.Amount); err != nil {
				logger.Error("failed to publish ledger updated event", "error", err)
			}
		case errEv := <-errChan:
			logger.Error("error", "transactionID", errEv.Data.TransactionID)
		case <-done:
			logger.Debug("shutting down")
			return nil
		}
	}
}

// transactionCreatedMsgHandlerFn returns a nats.MsgHandler that handles transaction.created messages.
// Successful ledger updates are sent to the ledgerChan.
// Failed ledger updates are sent to the errChan.
func transactionCreatedMsgHandlerFn(
	logger *slog.Logger,
	handler *handler.UpdateLedgerHandler,
	ledgerChan chan<- *handler.LedgerUpdateResult,
	errChan chan<- event.TransactionCreatedEvent,
) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var ev event.TransactionCreatedEvent
		err := json.Unmarshal(msg.Data, &ev)
		if err != nil {
			logger.Error("failed to unmarshal event", "error", err)
			return
		}

		results, err := handler.Handle(ev)
		if err != nil {
			logger.Error("failed to handle message", "error", err)
			errChan <- ev
			return
		}

		for _, r := range results {
			ledgerChan <- r
		}
	}
}
