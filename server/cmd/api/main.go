package main

import (
	"beerbux/internal/api"
	"beerbux/internal/api/middleware"
	"beerbux/internal/auth/command"
	authQueries "beerbux/internal/auth/db"
	authHandler "beerbux/internal/auth/handler"
	sessionHandler "beerbux/internal/session/handler"
	transactionHandler "beerbux/internal/transaction/handler"
	userHandler "beerbux/internal/user/handler"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log/slog"
	"net/http"
	"os"
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

	mux := buildServerMux(cfg, logger, database)
	logger.Debug("Starting server", "addr", cfg.APIAddress)
	if err := http.ListenAndServe(cfg.APIAddress, mux); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

type HandlerBuilderFunc func(cfg *api.Config, logger *slog.Logger, database *sql.DB, mux *http.ServeMux)

func buildServerMux(cfg *api.Config, logger *slog.Logger, database *sql.DB) http.Handler {
	mux := http.NewServeMux()

	builders := []HandlerBuilderFunc{
		authHandler.MakeHandlerRoutes,
		sessionHandler.MakeHandlerRoutes,
		transactionHandler.MakeHandlerRoutes,
		userHandler.MakeHandlerRoutes,
	}

	for _, buildRoutes := range builders {
		buildRoutes(cfg, logger, database, mux)
	}

	aq := authQueries.New(database)
	refreshTokenCommand := command.NewRefreshTokenCommand(aq, cfg.GetAuthOptions())
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
