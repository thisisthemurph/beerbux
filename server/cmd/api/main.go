package main

import (
	"beerbux/internal/api"
	"beerbux/internal/api/middleware"
	"beerbux/internal/auth/command"
	authQueries "beerbux/internal/auth/db"
	authHandler "beerbux/internal/auth/handler"
	sessionHandler "beerbux/internal/session/handler"
	"beerbux/internal/sse"
	streamHandler "beerbux/internal/streamer/handler"
	transactionHandler "beerbux/internal/transaction/handler"
	userHandler "beerbux/internal/user/handler"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := api.LoadConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	if err := run(logger, cfg); err != nil {
		logger.Error("Failed to start gateway-api", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, cfg *api.Config) error {
	logger.Debug("Starting API", "environment", cfg.Environment)

	database, err := connectToDatabase(cfg.Database.Driver, cfg.Database.URI)
	if err != nil {
		return err
	}

	streamServer := sse.NewServer(logger)
	ctx, cancel := createNotifyContext()
	defer cancel()

	// TODO: This is where the kafka consumer would add messages to the server
	msgChan := make(chan *sse.Message, 10)

	errChan := make(chan error, 1)
	mux := buildServerMux(cfg, logger, database, streamServer, msgChan)
	go func() {
		errChan <- http.ListenAndServe(cfg.APIAddress, mux)
	}()

	hb := time.NewTicker(time.Duration(cfg.StreamService.HeartbeatTickerSeconds) * time.Second)
	logger.Debug("Starting API server", "addr", cfg.APIAddress)

	for {
		select {
		case <-hb.C:
			streamServer.Heartbeat()
		case msg := <-msgChan:
			switch msg.Topic {
			case "session.transaction.created":
				streamServer.BroadcastMessageToRoom(msg.Key, msg)
			default:
				logger.Error("Unknown message topic", "topic", msg.Topic)
			}
		case err := <-errChan:
			return err
		case <-ctx.Done():
			logger.Debug("Shutting down API server")
			return nil
		}
	}
}

type HandlerBuilderFunc func(
	cfg *api.Config,
	logger *slog.Logger,
	database *sql.DB,
	mux *http.ServeMux,
	msgChan chan<- *sse.Message)

func buildServerMux(
	cfg *api.Config,
	logger *slog.Logger,
	database *sql.DB,
	streamServer *sse.Server,
	msgChan chan<- *sse.Message,
) http.Handler {
	mux := http.NewServeMux()

	builders := []HandlerBuilderFunc{
		authHandler.MakeHandlerRoutes,
		sessionHandler.MakeHandlerRoutes,
		transactionHandler.MakeHandlerRoutes,
		userHandler.MakeHandlerRoutes,
	}

	for _, buildRoutes := range builders {
		buildRoutes(cfg, logger, database, mux, msgChan)
	}

	mux.Handle("/events/session", streamHandler.NewSessionTransactionCreatedHandler(logger, streamServer))

	authenticationQueries := authQueries.New(database)
	refreshTokenCommand := command.NewRefreshTokenCommand(authenticationQueries, cfg.GetAuthOptions())
	authMiddleware := middleware.NewAuthMiddleware(refreshTokenCommand, cfg.Secrets.JWTSecret)
	recoverMiddleware := middleware.NewRecoverMiddleware(logger)
	return recoverMiddleware.Recover(authMiddleware.WithJWT(middleware.CORS(mux, cfg.ClientBaseURL)))
}

func connectToDatabase(driver, uri string) (*sql.DB, error) {
	database, err := sql.Open(driver, uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := migrateDatabase(database, driver); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return database, nil
}

func migrateDatabase(db *sql.DB, driver string) error {
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
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
