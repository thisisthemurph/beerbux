package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/thisisthemurph/beerbux/session-service/internal/application"
	"github.com/thisisthemurph/beerbux/session-service/internal/config"
	"github.com/thisisthemurph/beerbux/session-service/internal/handler"
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

	// Listen for user.updated events

	msgHandler := handler.NewUserUpdatedEventHandler(app.SessionRepository, app.Logger)
	_, err = app.NatsConn.Subscribe("user.updated", msgHandler.Handle)
	if err != nil {
		return fmt.Errorf("failed to subscribe to user.updated event: %w", err)
	}

	// gRPC Server

	grpcServer := app.NewSessionGRPCServer()

	listener, err := net.Listen("tcp", cfg.SessionServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
