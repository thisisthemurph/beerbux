package main

import (
	"database/sql"
	"fmt"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"log/slog"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/config"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
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

	w := repository.NewLedgerQueries(db)
	transactionCreatedMsgHandler := handler.NewTransactionCreatedMsgHandler(w, logger)
	_, err = nc.Subscribe("transaction.created", transactionCreatedMsgHandler.Handle)

	select {}
}
