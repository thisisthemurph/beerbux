package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"
	"github.com/thisisthemurph/beerbux/user-service/internal/config"
	"github.com/thisisthemurph/beerbux/user-service/internal/handler"
	"github.com/thisisthemurph/beerbux/user-service/internal/kafka"
	"github.com/thisisthemurph/beerbux/user-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/internal/server"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	if err := run(logger, cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, cfg *config.Config) error {
	logger.Info("Starting user-service", "env", cfg.Environment)
	if err := kafka.EnsureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return err
	}

	db, err := connectToDatabase(cfg.Database)
	if err != nil {
		return err
	}

	ctx, cancel := setupSignalHandler()
	defer cancel()

	setupAndRunKafkaConsumers(ctx, cfg, logger, db)

	listener, grpcServer, err := setupGRPCServer(logger, db, cfg.UserServerAddress, cfg.Kafka.Brokers, cfg.Environment)
	if err != nil {
		return fmt.Errorf("failed to setup gRPC server: %w", err)
	}
	defer grpcServer.GracefulStop()
	defer listener.Close()

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- grpcServer.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down")
		return nil
	case err := <-serverErrChan:
		return err
	}
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
) {
	newConsumer := func(topic string) kafka.ConsumerListener {
		return kafka.NewConsumer(logger, cfg.Kafka.Brokers, topic, "user-service")
	}

	userRepository := user.New(db)
	consumerHandlerMap := map[kafka.ConsumerListener]handler.KafkaMessageHandler{
		newConsumer("auth.user.registered"):          handler.NewAuthUserRegisteredHandler(userRepository),
		newConsumer("ledger.user.totals.calculated"): handler.NewLedgerUserTotalsCalculatedEventHandler(userRepository),
	}

	for c, h := range consumerHandlerMap {
		go c.StartListening(ctx, h.Handle)
	}
}

func setupGRPCServer(logger *slog.Logger, db *sql.DB, address string, brokers []string, environment config.Environment) (net.Listener, *grpc.Server, error) {
	userCreatedPublisher := publisher.NewUserCreatedKafkaPublisher(brokers)
	userUpdatedPublisher := publisher.NewUserUpdatedKafkaPublisher(brokers)

	userRepo := user.New(db)
	userServer := server.NewUserServer(userRepo, userCreatedPublisher, userUpdatedPublisher, logger)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServer(grpcServer, userServer)

	if environment.IsDevelopment() {
		reflection.Register(grpcServer)
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen at %v: %w", address, err)
	}

	return listener, grpcServer, nil
}

func connectToDatabase(dbCfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open(dbCfg.Driver, dbCfg.URI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrateDatabase(db, dbCfg.Driver); err != nil {
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
