package handler

import (
	"beerbux/internal/common/history"
	"beerbux/internal/common/sessionaccess"
	sessionaccessQueries "beerbux/internal/common/sessionaccess/db"
	"beerbux/internal/common/useraccess"
	useraccessQueries "beerbux/internal/common/useraccess/db"
	"beerbux/internal/session/command"
	"beerbux/internal/session/db"
	"beerbux/internal/session/query"
	"beerbux/internal/sse"
	"database/sql"
	"log/slog"
	"net/http"
)

func BuildRoutes(logger *slog.Logger, database *sql.DB, mux *http.ServeMux, msgChan chan<- *sse.Message) {
	queries := db.New(database)
	sessionHistoryService := history.NewSessionHistoryService(queries, logger)
	userReaderService := useraccess.NewUserReaderService(useraccessQueries.New(database))
	sessionReaderService := sessionaccess.NewSessionService(sessionaccessQueries.New(database))

	listSessionsByUserIDQuery := query.NewListSessionsByUserIDQuery(queries)
	createSessionCommand := command.NewCreateSessionCommand(database, queries, userReaderService)
	addSessionMemberCommand := command.NewAddSessionMemberCommand(database, queries, sessionHistoryService)
	removeSessionMemberCommand := command.NewRemoveSessionMemberCommand(queries, sessionHistoryService)
	updateSessionMemberAdminStateCommand := command.NewUpdateSessionMemberAdminStateCommand(queries, sessionHistoryService)
	updateSessionActiveStateCommand := command.NewUpdateSessionActionStateCommand(queries, sessionHistoryService)
	createTransactionCommand := command.NewCreateTransactionCommand(database, queries, sessionHistoryService)

	mux.Handle("GET /user/sessions", NewListCurrentUserSessionsHandler(listSessionsByUserIDQuery, logger))

	mux.Handle("GET /session/{sessionId}", NewGetSessionHandler(sessionReaderService, logger))
	mux.Handle("POST /session", NewCreateSessionHandler(createSessionCommand, logger))
	mux.Handle("POST /session/{sessionId}/member", NewAddSessionMemberHandler(userReaderService, sessionReaderService, addSessionMemberCommand, logger))
	mux.Handle("POST /session/{sessionId}/member/{memberId}/admin", NewUpdateSessionMemberAdminHandler(sessionReaderService, updateSessionMemberAdminStateCommand, logger))
	mux.Handle("DELETE /session/{sessionId}/member/{memberId}", NewRemoveSessionMemberHandler(sessionReaderService, removeSessionMemberCommand, logger))
	mux.Handle("DELETE /session/{sessionId}/leave", NewLeaveSessionHandler(removeSessionMemberCommand, logger))
	mux.Handle("PUT /session/{sessionId}/state/{command}", NewUpdateSessionActiveStateHandler(sessionReaderService, updateSessionActiveStateCommand, logger))

	mux.Handle("GET /session/{sessionId}/history", NewGetSessionHistoryHandler(sessionReaderService, sessionHistoryService, logger))

	mux.Handle("POST /session/{sessionId}/transaction", NewCreateTransactionHandler(sessionReaderService, createTransactionCommand, logger, msgChan))
}
