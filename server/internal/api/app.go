package api

import (
	"beerbux/internal/api/config"
	"beerbux/internal/api/database"
	"beerbux/internal/sse"
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Config      *config.Config
	Logger      *slog.Logger
	DB          *sql.DB
	messageChan chan *sse.Message
}

func NewApp(cfg *config.Config, logger *slog.Logger) (*App, error) {
	logger.Debug("Starting API", "environment", cfg.Environment)

	db, err := database.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:      cfg,
		Logger:      logger,
		DB:          db,
		messageChan: make(chan *sse.Message, 10),
	}, nil
}

func (app *App) MessageReceiver() chan<- *sse.Message {
	return app.messageChan
}

func (app *App) Start() error {
	streamServer := sse.NewServer(app.Logger)
	ctx, cancel := createNotifyContext()
	defer cancel()

	errChan := make(chan error, 1)
	server, err := app.NewServer(streamServer)
	if err != nil {
		app.Logger.Error("Failed to build routes", "error", err)
		return err
	}

	go func() {
		errChan <- server.ListenAndServe()
	}()

	hb := time.NewTicker(time.Duration(app.Config.StreamService.HeartbeatTickerSeconds) * time.Second)
	app.Logger.Debug("Starting API server", "addr", app.Config.Address)

	for {
		select {
		case <-hb.C:
			streamServer.Heartbeat()
		case msg := <-app.messageChan:
			switch msg.Topic {
			case "session.transaction.created":
				streamServer.BroadcastMessageToRoom(msg.Key, msg)
			default:
				app.Logger.Error("Unknown message topic", "topic", msg.Topic)
			}
		case err := <-errChan:
			app.Logger.Error("Error during API server", "error", err)
			return err
		case <-ctx.Done():
			app.Logger.Debug("Shutting down API server")
			return nil
		}
	}
}

func createNotifyContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	return ctx, cancel
}
