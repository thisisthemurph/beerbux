package handler

import (
	"beerbux/internal/api/config"
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/db"
	"beerbux/pkg/email"
	"database/sql"
	"log/slog"
	"net/http"
)

func BuildRoutes(
	cfg *config.Config,
	logger *slog.Logger,
	database *sql.DB,
	emailSender email.Sender,
	mux *http.ServeMux,
) {
	options := cfg.GetAuthOptions()

	queries := db.New(database)
	loginCommand := command.NewLoginCommand(queries, options)
	signupCommand := command.NewSignupCommand(queries, options)
	refreshCommand := command.NewRefreshTokenCommand(queries, options)
	invalidateRefreshTokenCommand := command.NewInvalidateRefreshTokenCommand(queries)
	initializePasswordResetCommand := command.NewInitializePasswordResetCommand(queries)
	resetPasswordCommand := command.NewResetPasswordCommand(queries)

	mux.Handle("POST /auth/login", NewLoginHandler(loginCommand, logger))
	mux.Handle("POST /auth/signup", NewSignupHandler(signupCommand, logger))
	mux.Handle("POST /auth/refresh", NewRefreshHandler(refreshCommand, logger))
	mux.Handle("POST /auth/logout", NewLogoutHandler(invalidateRefreshTokenCommand, logger))
	mux.Handle("POST /auth/password/initialize-reset", NewInitializePasswordUpdateHandler(
		initializePasswordResetCommand, emailSender, logger))
	mux.Handle("POST /auth/password/reset", NewUpdatePasswordHandler(resetPasswordCommand, logger))
}
