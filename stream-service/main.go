package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thisisthemurph/beerbux/stream-service/internal/config"
	"github.com/thisisthemurph/beerbux/stream-service/internal/consumer"
	"github.com/thisisthemurph/beerbux/stream-service/internal/handler"
	"github.com/thisisthemurph/beerbux/stream-service/internal/middleware"
	"github.com/thisisthemurph/beerbux/stream-service/internal/sse"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	if err := run(cfg, logger); err != nil {
		logger.Error("Error running server", "error", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("Starting stream-service", "env", cfg.Environment)

	server := sse.NewServer(logger)

	ctx, cancel := setupSignalHandler()
	defer cancel()

	setupAndRunKafkaConsumers(ctx, logger, cfg.Kafka.Brokers, server)

	serverErrChan := make(chan error, 1)
	setupAndRunServerSentEventsServer(cfg, logger, server, serverErrChan)

	hb := time.NewTicker(time.Duration(cfg.HeartbeatTickerSeconds) * time.Second)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutting down server...")
			return nil
		case <-hb.C:
			server.Heartbeat()
		case err := <-serverErrChan:
			logger.Error("Error starting server:", "error", err)
			return err
		}
	}
}

func setupAndRunServerSentEventsServer(cfg *config.Config, logger *slog.Logger, server *sse.Server, errCh chan<- error) {
	logger.Info("Starting server-sent events server")
	mux := setupServerMux(logger, cfg.ClientBaseURL, server)

	go func() {
		errCh <- http.ListenAndServe(cfg.StreamServerAddr, mux)
	}()
}

func setupServerMux(logger *slog.Logger, clientBaseURL string, server *sse.Server) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/events/session", handler.NewSessionTransactionCreatedHandler(logger, server))

	return middleware.CORS(mux, clientBaseURL)
}

func setupAndRunKafkaConsumers(ctx context.Context, logger *slog.Logger, brokers []string, server *sse.Server) {
	logger.Info("Setting up Kafka consumers")
	transactionCreatedConsumer := consumer.NewConsumer(logger, brokers, "session.transaction.created", "stream-service")
	go transactionCreatedConsumer.StartListening(ctx, consumer.NewSessionTransactionCreatedKafkaConsumer(server))
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
