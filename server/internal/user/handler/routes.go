package handler

import (
	"beerbux/internal/common/useraccess"
	useraccessQueries "beerbux/internal/common/useraccess/db"
	"database/sql"
	"log/slog"
	"net/http"
)

func BuildRoutes(
	logger *slog.Logger,
	database *sql.DB,
	mux *http.ServeMux,
) {
	uaQueries := useraccessQueries.New(database)
	userReaderService := useraccess.NewUserReaderService(uaQueries)

	mux.Handle("GET /api/user", NewGetCurrentUserHandler(userReaderService, logger))
	mux.Handle("GET /api/user/{userId}/balance", NewGetCurrentUserBalanceHandler(userReaderService, logger))
}
