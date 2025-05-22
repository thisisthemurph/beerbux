package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/history"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"log/slog"
	"net/http"
)

type GetSessionHistoryHandler struct {
	sessionHistoryReader history.SessionHistoryReader
	getSessionQuery      *query.GetSessionQuery
	logger               *slog.Logger
}

func NewGetSessionHistoryHandler(sessionHistoryReader history.SessionHistoryReader, getSessionQuery *query.GetSessionQuery, logger *slog.Logger) *GetSessionHistoryHandler {
	return &GetSessionHistoryHandler{
		sessionHistoryReader: sessionHistoryReader,
		getSessionQuery:      getSessionQuery,
		logger:               logger,
	}
}

func (h *GetSessionHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, ok := url.Path.GetUUID(r, "sessionId")
	if !ok {
		send.BadRequest(w, "Session ID required")
		return
	}

	s, err := h.getSessionQuery.Execute(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "The session could not be found")
			return
		}
		send.InternalServerError(w, "There was an issue finding the session")
		return
	}

	if !s.IsMember(c.Subject) {
		send.Unauthorized(w, "You are not a member of this session")
		return
	}

	hist, err := h.sessionHistoryReader.GetSessionHistory(r.Context(), sessionID)
	if err != nil {
		send.InternalServerError(w, "There has been an error fetching the session history")
		return
	}

	send.JSON(w, hist, http.StatusOK)
}
