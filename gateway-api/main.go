package main

import (
	"fmt"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/session"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/transaction"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/transaction-service/protos/transactionpb"
	"log/slog"
	"net/http"
	"os"

	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/config"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/auth"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/user"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/middleware"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
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

	userClientConn, err := grpc.NewClient(cfg.UserServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error connecting to user server: %w", err)
	}
	userClient := userpb.NewUserClient(userClientConn)

	sessionClientConn, err := grpc.NewClient(cfg.SessionServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error connecting to session server: %w", err)
	}
	sessionClient := sessionpb.NewSessionClient(sessionClientConn)

	transactionConn, err := grpc.NewClient(cfg.TransactionServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error connecting to transaction server: %w", err)
	}
	transactionClient := transactionpb.NewTransactionClient(transactionConn)

	mux := buildServerMux(cfg, authClient, userClient, sessionClient, transactionClient)
	logger.Debug("Starting server", "add", cfg.GatewayAPIAddress)
	if err := http.ListenAndServe(cfg.GatewayAPIAddress, mux); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func buildServerMux(
	cfg *config.Config,
	authClient authpb.AuthClient,
	userClient userpb.UserClient,
	sessionClient sessionpb.SessionClient,
	transactionClient transactionpb.TransactionClient,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /api/auth/login", auth.NewLoginHandler(authClient, userClient))
	mux.Handle("POST /api/auth/signup", auth.NewSignupHandler(authClient))
	mux.Handle("POST /api/auth/refresh", auth.NewRefreshHandler(authClient))
	mux.Handle("POST /api/auth/logout", auth.NewLogoutHandler(authClient))

	mux.Handle("GET /api/user", user.NewGetCurrentUserHandler(userClient))
	mux.Handle("GET /api/user/sessions", session.NewListSessionsForUserHandler(sessionClient))

	mux.Handle("GET /api/session/{sessionId}", session.NewGetSessionByIdHandler(sessionClient))
	mux.Handle("POST /api/session", session.NewCreateSessionHandler(sessionClient))
	mux.Handle("POST /api/session/{sessionId}/member", session.NewAddMemberToSessionHandler(userClient, sessionClient))
	mux.Handle("POST /api/session/{sessionId}/transaction", transaction.NewCreateTransactionHandler(transactionClient))

	authMiddleware := middleware.NewAuthMiddleware(authClient, cfg.Secrets.JWTSecret)
	return middleware.Recover(authMiddleware.WithJWT(middleware.CORS(mux, cfg.ClientBaseURL)))
}
