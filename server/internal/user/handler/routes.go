package handler

import (
	"beerbux/internal/common/useraccess"
	useraccessQueries "beerbux/internal/common/useraccess/db"
	"beerbux/internal/user/command"
	"beerbux/internal/user/db"
	"database/sql"
	"log/slog"
	"net/http"
)

func BuildRoutes(
	logger *slog.Logger,
	database *sql.DB,
	mux *http.ServeMux,
) {
	queries := db.New(database)
	uaQueries := useraccessQueries.New(database)
	userReaderService := useraccess.NewUserReaderService(uaQueries)

	updateUserCommand := command.NewUpdateUserCommand(queries)

	mux.Handle("GET /user", NewGetCurrentUserHandler(userReaderService, logger))
	mux.Handle("PUT /user", NewUpdateUserHandler(updateUserCommand, userReaderService, logger))
	mux.Handle("GET /user/{userId}/balance", NewGetCurrentUserBalanceHandler(userReaderService, logger))
}
