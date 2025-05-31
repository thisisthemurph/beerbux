package handler

import (
	"beerbux/internal/api/config"
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/db"
	"beerbux/internal/common/useraccess"
	useraccessQueries "beerbux/internal/common/useraccess/db"
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
	generateTokensCommand := command.NewGenerateTokensCommand(queries, options)
	signupCommand := command.NewSignupCommand(queries, options)
	refreshCommand := command.NewRefreshTokenCommand(queries, options)
	invalidateRefreshTokenCommand := command.NewInvalidateRefreshTokenCommand(queries)
	initializeUpdatePasswordCommand := command.NewInitializeUpdatePasswordCommand(queries)
	updatePasswordCommand := command.NewUpdatePasswordCommand(queries)
	initializeUpdateEmailCommand := command.NewInitializeUpdateEmailCommand(queries)
	updateEmailCommand := command.NewUpdateEmailCommand(queries)
	comparePasswordCommand := command.NewComparePasswordCommand(queries)
	initializePasswordResetCommand := command.NewInitializePasswordResetCommand(queries)
	resetPasswordCommand := command.NewResetPasswordCommand(queries)

	userAccessQueries := useraccessQueries.New(database)
	userReaderService := useraccess.NewUserReaderService(userAccessQueries)

	mux.Handle("POST /auth/login", NewLoginHandler(generateTokensCommand, comparePasswordCommand, logger))
	mux.Handle("POST /auth/signup", NewSignupHandler(signupCommand, logger))
	mux.Handle("POST /auth/refresh", NewRefreshHandler(refreshCommand, logger))
	mux.Handle("POST /auth/logout", NewLogoutHandler(invalidateRefreshTokenCommand, logger))
	mux.Handle("POST /auth/password/initialize-update", NewInitializeUpdatePasswordHandler(initializeUpdatePasswordCommand, emailSender, logger))
	mux.Handle("PUT /auth/password", NewUpdatePasswordHandler(updatePasswordCommand, logger))
	mux.Handle("POST /auth/password/initialize-reset", NewInitializePasswordResetHandler(initializePasswordResetCommand, userReaderService, emailSender, logger))
	mux.Handle("PUT /auth/password/reset", NewResetPasswordHandler(resetPasswordCommand, logger))
	mux.Handle("POST /auth/email/initialize-update", NewInitializeEmailUpdateHandler(initializeUpdateEmailCommand, userReaderService, emailSender, logger))
	mux.Handle("PUT /auth/email", NewUpdateEmailHandler(updateEmailCommand, generateTokensCommand, logger))
}
