package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/auth-service/internal/config"
	"github.com/thisisthemurph/beerbux/auth-service/internal/producer"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/auth"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/token"
	"github.com/thisisthemurph/beerbux/auth-service/internal/server"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       cfg.SlogLevel(),
		ReplaceAttr: nil,
	}))

	if err := run(logger, cfg); err != nil {
		logger.Error("run error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, cfg *config.Config) error {
	logger.Debug("connecting to database")
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := ensureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return fmt.Errorf("failed to ensure Kafka topics: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	serverErrChan := make(chan error, 1)

	// Auth server
	listener, grpcServer, err := makeAuthServer(logger, db, cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Error("failed to close listener", "error", err)
		}
		grpcServer.GracefulStop()
	}()

	go func() {
		logger.Debug("listening", "address", cfg.AuthServerAddress)
		serverErrChan <- grpcServer.Serve(listener)
	}()

	select {
	case <-sigChan:
		logger.Debug("shutting down")
		return nil
	case err := <-serverErrChan:
		return fmt.Errorf("server error: %w", err)
	}
}

func makeAuthServer(logger *slog.Logger, db *sql.DB, cfg *config.Config) (net.Listener, *grpc.Server, error) {
	logger.Debug("starting server")
	grpcServer := grpc.NewServer()
	authRepo := auth.New(db)
	authTokenRepo := token.New(db)

	userRegisteredProducer := producer.NewUserRegisteredKafkaProducer(cfg.Kafka.Brokers)

	authServer := server.NewAuthServer(logger, authRepo, authTokenRepo, userRegisteredProducer, server.AuthServerOptions{
		JWTSecret:       cfg.Secrets.JWTSecret,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	})
	authpb.RegisterAuthServer(grpcServer, authServer)

	if cfg.Environment.IsDevelopment() {
		reflection.Register(grpcServer)
	}

	listener, err := net.Listen("tcp", cfg.AuthServerAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	return listener, grpcServer, nil
}

func ensureKafkaTopics(brokers []string) error {
	if err := kafkatopic.EnsureTopicExists(brokers, kafka.TopicConfig{
		Topic:             producer.TopicUserRegistered,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		return fmt.Errorf("failed to ensure session.member.added Kafka topic: %w", err)
	}

	return nil
}
