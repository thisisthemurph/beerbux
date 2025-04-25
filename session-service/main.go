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
	"github.com/thisisthemurph/beerbux/session-service/internal/config"
	"github.com/thisisthemurph/beerbux/session-service/internal/handler"
	"github.com/thisisthemurph/beerbux/session-service/internal/kafka"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/history"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/internal/server"
	"github.com/thisisthemurph/beerbux/session-service/protos/historypb"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	if err := run(cfg, logger); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("Starting session-service", "env", cfg.Environment)

	if err := kafka.EnsureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return fmt.Errorf("failed to ensure Kafka topics: %w", err)
	}

	db, err := connectToDatabase(cfg.Database)
	if err != nil {
		return err
	}

	userClientConn, err := grpc.NewClient(cfg.UserServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error connecting to user-service: %w", err)
	}
	defer userClientConn.Close()

	ctx, cancel := setupSignalHandler()
	defer cancel()

	sessionRepository := session.New(db)
	historyRepository := history.NewHistoryRepository(db)

	// Setup and run Kafka consumers.
	setupAndRunKafkaConsumers(ctx, cfg, logger, db, sessionRepository)

	// Set up the session gRPC server.
	listener, grpcServer, err := setupSessionGRPCServer(
		cfg,
		logger,
		db,
		sessionRepository,
		historyRepository,
		userClientConn)
	if err != nil {
		return err
	}

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- grpcServer.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down gRPC server")
		grpcServer.GracefulStop()
		listener.Close()
	case err := <-serverErrChan:
		// The GRPC server has had an issue listening, and we should quit the app.
		return fmt.Errorf("gRPC server error: %w", err)
	}

	return nil
}

func connectToDatabase(dbCfg config.DBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("file:%s?_busy_timeout=3000&cache=shared&_fk=true", dbCfg.URI)
	db, err := sql.Open(dbCfg.Driver, connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
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
	sessionRepository *session.Queries,
) {
	newConsumer := func(topic string) kafka.ConsumerListener {
		return kafka.NewConsumer(logger, cfg.Kafka.Brokers, topic, "session-service")
	}

	sessionRepositoryWithTransactions := repository.NewSessionQueries(db)
	historyRepository := history.NewHistoryRepository(db)
	sessionTransactionCreatedPublisher := publisher.NewSessionTransactionCreatedKafkaPublisher(cfg.Kafka.Brokers)
	consumerHandlerMap := map[kafka.ConsumerListener]handler.KafkaMessageHandler{
		newConsumer("user.updated"):        handler.NewUserUpdatedEventHandler(sessionRepository),
		newConsumer("transaction.created"): handler.NewTransactionCreatedMessageHandler(sessionRepositoryWithTransactions, historyRepository, sessionTransactionCreatedPublisher),
	}

	for consumer, h := range consumerHandlerMap {
		go consumer.StartListening(ctx, h.Handle)
	}
}

func setupSessionGRPCServer(
	cfg *config.Config,
	logger *slog.Logger,
	db *sql.DB,
	sessionRepository *session.Queries,
	historyRepository history.HistoryRepository,
	userClientConn *grpc.ClientConn,
) (net.Listener, *grpc.Server, error) {
	listener, err := net.Listen("tcp", cfg.SessionServerAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	ss := server.NewSessionServer(
		db,
		sessionRepository,
		historyRepository,
		userpb.NewUserClient(userClientConn),
		publisher.NewSessionMemberAddedKafkaPublisher(cfg.Kafka.Brokers),
		logger)

	hs := server.NewHistoryServer(historyRepository)

	sessionpb.RegisterSessionServer(grpcServer, ss)
	historypb.RegisterHistoryServer(grpcServer, hs)

	if cfg.Environment.IsDevelopment() {
		reflection.Register(grpcServer)
	}

	return listener, grpcServer, nil
}
