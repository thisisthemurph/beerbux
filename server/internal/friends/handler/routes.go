package handler

import (
	"beerbux/internal/friends/db"
	"beerbux/internal/friends/query"
	"database/sql"
	"log/slog"
	"net/http"
)

func BuildRoutes(logger *slog.Logger, database *sql.DB, muc *http.ServeMux) {
	queries := db.New(database)
	getFriendsQuery := query.NewGetFriendsQuery(queries)

	muc.Handle("GET /friends", NewGetFriendsHandler(getFriendsQuery, logger))
}
