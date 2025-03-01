package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/nats-io/nats.go"
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

	logger.Debug("connecting to the database")
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return err
	}
	defer db.Close()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	userCreatedPublisher := publisher.NewUserCreatedNatsPublisher(nc)
	userUpdatedPublisher := publisher.NewUserUpdatedNatsPublisher(nc)

	logger.Debug("connected to NATS", "url", nats.DefaultURL)

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
