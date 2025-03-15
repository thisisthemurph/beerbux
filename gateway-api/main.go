package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/config"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/auth"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLevel(),
	}))

	if err := run(logger, cfg); err != nil {
		logger.Error("Failed to start gateway-api", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, cfg *config.Config) error {
	logger.Debug("Starting gateway-api", "environment", cfg.Environment)

	authClientConn, err := grpc.NewClient(cfg.AuthServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error connecting to auth server: %w", err)
	}
	authClient := authpb.NewAuthClient(authClientConn)

	mux := buildServerMux(cfg, authClient)
	logger.Debug("Starting server", "add", cfg.GatewayAPIAddress)
	if err := http.ListenAndServe(cfg.GatewayAPIAddress, mux); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func buildServerMux(cfg *config.Config, authClient authpb.AuthClient) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /auth/login", auth.NewLoginHandler(authClient))
	mux.Handle("POST /auth/signup", auth.NewSignupHandler(authClient))
	mux.Handle("POST /auth/refresh", auth.NewRefreshHandler(authClient))

	return middleware.Recover(middleware.WithJWT(middleware.CORS(mux, cfg.ClientBaseURL), cfg.Secrets.JWTSecret))
}
