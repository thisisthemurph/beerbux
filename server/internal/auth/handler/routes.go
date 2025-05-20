package handler

import (
	"beerbux/internal/api"
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/db"
	"beerbux/internal/sse"
	"database/sql"
	"log/slog"
	"net/http"
)

func MakeHandlerRoutes(
	cfg *api.Config,
	logger *slog.Logger,
	database *sql.DB,
	mux *http.ServeMux,
	_ chan<- *sse.Message,
) {
	options := cfg.GetAuthOptions()

	queries := db.New(database)
	loginCommand := command.NewLoginCommand(queries, options)
	signupCommand := command.NewSignupCommand(queries, options)
	refreshCommand := command.NewRefreshTokenCommand(queries, options)
	invalidateRefreshTokenCommand := command.NewInvalidateRefreshTokenCommand(queries)

	mux.Handle("POST /api/auth/login", NewLoginHandler(loginCommand, logger))
	mux.Handle("POST /api/auth/signup", NewSignupHandler(signupCommand, logger))
	mux.Handle("POST /api/auth/refresh", NewRefreshHandler(refreshCommand, logger))
	mux.Handle("POST /api/auth/logout", NewLogoutHandler(invalidateRefreshTokenCommand, logger))
}
