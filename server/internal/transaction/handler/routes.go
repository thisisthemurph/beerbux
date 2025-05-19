package handler

import (
	"beerbux/internal/api"
	"beerbux/internal/common/history"
	sessionDB "beerbux/internal/session/db"
	sessionQuery "beerbux/internal/session/query"
	"beerbux/internal/transaction/command"
	"beerbux/internal/transaction/db"
	"database/sql"
	"log/slog"
	"net/http"
)

func MakeHandlerRoutes(_ *api.Config, logger *slog.Logger, database *sql.DB, mux *http.ServeMux) {
	queries := db.New(database)
	ssnQueries := sessionDB.New(database)
	sessionHistoryService := history.NewSessionHistoryService(ssnQueries, logger)

	getSessionQuery := sessionQuery.NewGetSessionQuery(ssnQueries)
	createTransactionCommand := command.NewCreateTransactionCommand(database, queries, sessionHistoryService)

	mux.Handle("POST /api/session/{sessionId}/transaction", NewCreateTransactionHandler(getSessionQuery, createTransactionCommand, logger))
}
