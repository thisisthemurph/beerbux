package main

import (
	"context"
	"fmt"
	"github.com/thisisthemurph/beerbux/session-service/internal/handler"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/session-service/internal/application"
	"github.com/thisisthemurph/beerbux/session-service/internal/config"
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

	app, err := application.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	defer app.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Listen for user.updated Kafka events
	h := handler.NewUserUpdatedEventHandler(app.SessionRepository)
	go listenForUserUpdatedEvents(ctx, logger, cfg.Kafka.Brokers, h)

	// Start gRPC server
	grpcServer := app.NewSessionGRPCServer()
	listener, err := net.Listen("tcp", cfg.SessionServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- grpcServer.Serve(listener)
	}()

	select {
	case <-sigChan:
		logger.Info("shutting down")
	case err := <-serverErrChan:
		if err != nil {
			logger.Error("gRPC server error", "error", err)
		}
	}

	logger.Info("gracefully shutting down gRPC server")
	grpcServer.GracefulStop()

	cancel()
	logger.Info("gracefully shut down Kafka consumer")
	return nil
}

func listenForUserUpdatedEvents(ctx context.Context, logger *slog.Logger, brokers []string, h *handler.UserUpdatedEventHandler) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "user.updated",
		GroupID: "session-service",
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Kafka consumer shutting down")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				logger.Error("failed to read message", "error", err)
				continue
			}
			if err := h.Handle(ctx, msg); err != nil {
				logger.Error("failed to handle message", "error", err)
			}
		}
	}
}
