package application

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/session-service/internal/config"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/internal/server"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "modernc.org/sqlite"
)

type App struct {
	DB                *sql.DB
	NatsConn          *nats.Conn
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

	logger.Debug("connecting to NATS")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	logger.Debug("creating user-service gRPC client connection")
	userClientConn, err := grpc.NewClient(cfg.UserServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error connecting to user-service: %w", err)
	}

	return &App{
		DB:                db,
		NatsConn:          nc,
		SessionRepository: session.New(db),
		UserClient:        userpb.NewUserClient(userClientConn),
		Logger:            logger,

		cfg:                  cfg,
		userServerGRPCClient: userClientConn,
		built:                true,
	}, nil
}

func (app *App) Close() {
	app.userServerGRPCClient.Close()
	app.NatsConn.Close()
	app.DB.Close()
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
		app.Logger,
	)

	gs := grpc.NewServer()
	sessionpb.RegisterSessionServer(gs, ss)

	if app.cfg.Environment.IsDevelopment() {
		reflection.Register(gs)
	}

	return gs
}
