package application

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/session-service/internal/config"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/internal/server"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/shared/kafkatopic"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "modernc.org/sqlite"
)

type App struct {
	DB                *sql.DB
	SessionRepository *session.Queries
	UserClient        userpb.UserClient
	Logger            *slog.Logger

	cfg                  *config.Config
	userServerGRPCClient *grpc.ClientConn
	built                bool
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	logger.Debug("connecting to database")
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrateDatabase(db, cfg); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Debug("creating user-service gRPC client connection")
	userClientConn, err := grpc.NewClient(cfg.UserServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error connecting to user-service: %w", err)
	}

	if err := ensureKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return nil, fmt.Errorf("failed to ensure Kafka topics: %w", err)
	}

	return &App{
		DB:                db,
		SessionRepository: session.New(db),
		UserClient:        userpb.NewUserClient(userClientConn),
		Logger:            logger,

		cfg:                  cfg,
		userServerGRPCClient: userClientConn,
		built:                true,
	}, nil
}

func (app *App) Close() {
	_ = app.userServerGRPCClient.Close()
	_ = app.DB.Close()
}

func (app *App) NewSessionGRPCServer() *grpc.Server {
	app.Logger.Debug("creating session-service gRPC server")
	if !app.built {
		panic("app not built")
	}

	ss := server.NewSessionServer(
		app.DB,
		app.SessionRepository,
		app.UserClient,
		publisher.NewSessionMemberAddedKafkaPublisher(app.cfg.Kafka.Brokers),
		app.Logger,
	)

	gs := grpc.NewServer()
	sessionpb.RegisterSessionServer(gs, ss)

	if app.cfg.Environment.IsDevelopment() {
		reflection.Register(gs)
	}

	return gs
}

func ensureKafkaTopics(brokers []string) error {
	if err := kafkatopic.EnsureTopicExists(brokers, kafka.TopicConfig{
		Topic:             publisher.TopicSessionMemberAdded,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		return fmt.Errorf("failed to ensure session.member.added Kafka topic: %w", err)
	}

	return nil
}

func migrateDatabase(db *sql.DB, cfg *config.Config) error {
	if err := goose.SetDialect(cfg.Database.Driver); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Up(db, "./internal/db/migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
