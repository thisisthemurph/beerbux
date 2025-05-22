package routes

import (
	"beerbux/internal/api/config"
	"beerbux/internal/api/middleware"
	"beerbux/internal/auth/command"
	authQueries "beerbux/internal/auth/db"
	authHandler "beerbux/internal/auth/handler"
	sessionHandler "beerbux/internal/session/handler"
	"beerbux/internal/sse"
	streamHandler "beerbux/internal/streamer/handler"
	userHandler "beerbux/internal/user/handler"
	"database/sql"
	"log/slog"
	"net/http"
)

func Build(
	cfg *config.Config,
	logger *slog.Logger,
	database *sql.DB,
	streamServer *sse.Server,
	msgChan chan<- *sse.Message,
) http.Handler {
	mux := http.NewServeMux()

	authHandler.BuildRoutes(cfg, logger, database, mux)
	sessionHandler.BuildRoutes(logger, database, mux, msgChan)
	//transactionHandler.BuildRoutes(logger, database, mux, msgChan)
	userHandler.BuildRoutes(logger, database, mux)

	mux.Handle("/events/session", streamHandler.NewSessionTransactionCreatedHandler(logger, streamServer))

	authenticationQueries := authQueries.New(database)
	refreshTokenCommand := command.NewRefreshTokenCommand(authenticationQueries, cfg.GetAuthOptions())
	authMiddleware := middleware.NewAuthMiddleware(refreshTokenCommand, cfg.Secrets.JWTSecret)
	recoverMiddleware := middleware.NewRecoverMiddleware(logger)
	return recoverMiddleware.Recover(authMiddleware.WithJWT(middleware.CORS(mux, cfg.ClientBaseURL)))
}
