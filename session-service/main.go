package main

import (
	"fmt"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/internal/server"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"

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

	_, err = app.NatsConn.Subscribe("user.created", app.UserCreatedMsgHandler.Handle)
	if err != nil {
		return fmt.Errorf("failed to subscribe to user.created: %w", err)
	}

	sessionRepository := session.New(app.DB)
	sessionServer := server.NewSessionServer(sessionRepository, logger)

	grpcServer := grpc.NewServer()
	sessionpb.RegisterSessionServer(grpcServer, sessionServer)

	if cfg.Environment.IsDevelopment() {
		reflection.Register(grpcServer)
	}

	lis, err := net.Listen("tcp", cfg.SessionServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen at %v: %w", cfg.SessionServerAddress, err)
	}
	defer lis.Close()

	logger.Debug("starting session server", "address", cfg.SessionServerAddress)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
