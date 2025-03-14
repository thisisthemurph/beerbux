package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/thisisthemurph/beerbux/auth-service/internal/config"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/auth"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/token"
	"github.com/thisisthemurph/beerbux/auth-service/internal/server"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
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

	logger.Debug("starting server")
	grpcServer := grpc.NewServer()
	authRepo := auth.New(db)
	authTokenRepo := token.New(db)
	authServer := server.NewAuthServer(authRepo, authTokenRepo, server.AuthServerOptions{
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
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	logger.Debug("listening", "address", cfg.AuthServerAddress)
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
