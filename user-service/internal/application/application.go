package application

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/thisisthemurph/beerbux/user-service/internal/config"
	"github.com/thisisthemurph/beerbux/user-service/internal/consumer"
	"github.com/thisisthemurph/beerbux/user-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/ledger"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/internal/server"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	logger                *slog.Logger
	db                    *sql.DB
	grpcServer            *grpc.Server
	grpcListener          net.Listener
	ledgerUpdatedConsumer *consumer.LedgerUpdatedKafkaConsumer
}

func New(logger *slog.Logger, cfg *config.Config) (*App, error) {
	app := &App{logger: logger}

	if err := app.connectToDatabase(cfg.Database.Driver, cfg.Database.URI); err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	app.setupKafka(cfg.Kafka.Brokers)

	if err := app.setupGRPCServer(cfg.UserServerAddress, cfg.Kafka.Brokers, cfg.Environment); err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to setup gRPC server: %w", err)
	}

	return app, nil
}

func (a *App) RunLedgerUpdatedEventConsumer(ctx context.Context) {
	a.ledgerUpdatedConsumer.Listen(ctx)
}

func (a *App) RunGRPCServer() error {
	return a.grpcServer.Serve(a.grpcListener)
}

func (a *App) connectToDatabase(driver, uri string) error {
	a.logger.Debug("connecting to the database")
	db, err := sql.Open(driver, uri)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	a.db = db
	return nil
}

func (a *App) setupKafka(brokers []string) {
	a.logger.Debug("setting up Kafka consumers", "brokers", brokers)
	userLedgerRepo := ledger.New(a.db)
	ledgerUpdatedConsumer := consumer.NewLedgerUpdatedKafkaConsumer(
		a.logger, brokers, "ledger.updated", userLedgerRepo)

	a.ledgerUpdatedConsumer = ledgerUpdatedConsumer
}

func (a *App) setupGRPCServer(address string, brokers []string, environment config.Environment) error {
	a.logger.Debug("setting up Kafka publishers", "brokers", brokers)
	userCreatedPublisher := publisher.NewUserCreatedKafkaPublisher(brokers)
	userUpdatedPublisher := publisher.NewUserUpdatedKafkaPublisher(brokers)

	a.logger.Debug("setting up gRPC server", "address", address)
	userRepo := user.New(a.db)
	ledgerRepo := ledger.New(a.db)
	userServer := server.NewUserServer(userRepo, ledgerRepo, userCreatedPublisher, userUpdatedPublisher, a.logger)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServer(grpcServer, userServer)

	if environment.IsDevelopment() {
		reflection.Register(grpcServer)
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen at %v: %w", address, err)
	}

	a.grpcServer = grpcServer
	a.grpcListener = listener
	return nil
}

func (a *App) Close() {
	safeClose(a.logger, a.grpcListener)

	if a.grpcServer != nil {
		a.logger.Info("shutting down gRPC server")
		a.grpcServer.GracefulStop()
	}

	safeClose(a.logger, a.db)
	safeClose(a.logger, a.ledgerUpdatedConsumer)
}

func safeClose(logger *slog.Logger, closer io.Closer) {
	logger.Debug("closing resource", "type", fmt.Sprintf("%T", closer))
	if closer != nil {
		if err := closer.Close(); err != nil {
			logger.Error("error closing resource", "type", fmt.Sprintf("%T", closer), "error", err)
		}
	}
}
