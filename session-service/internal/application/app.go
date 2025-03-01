package application

import (
	"database/sql"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/session-service/internal/config"
	"github.com/thisisthemurph/beerbux/session-service/internal/handler"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/user"

	_ "modernc.org/sqlite"
)

type App struct {
	DB                    *sql.DB
	NatsConn              *nats.Conn
	UserCreatedMsgHandler *handler.UserCreatedMsgHandler
	Logger                *slog.Logger
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

	userRepo := user.New(db)
	userCreatedMgsHandler := handler.NewUserCreatedMsgHandler(userRepo)

	return &App{
		NatsConn:              nc,
		Logger:                logger,
		UserCreatedMsgHandler: userCreatedMgsHandler,
		DB:                    db,
	}, nil
}

func (app *App) Close() {
	app.NatsConn.Close()
	app.DB.Close()
}
