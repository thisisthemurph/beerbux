package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/history"
	"beerbux/internal/common/sessionaccess"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"log/slog"
	"net/http"
)

type GetSessionHistoryHandler struct {
	sessionReader        sessionaccess.SessionReader
	sessionHistoryReader history.SessionHistoryReader
	logger               *slog.Logger
}

func NewGetSessionHistoryHandler(sr sessionaccess.SessionReader, sessionHistoryReader history.SessionHistoryReader, logger *slog.Logger) *GetSessionHistoryHandler {
	return &GetSessionHistoryHandler{
		sessionReader:        sr,
		sessionHistoryReader: sessionHistoryReader,
		logger:               logger,
	}
}

// GetSessionHistoryHandler godoc
// @Summary Get Session History
// @Tags session
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} history.SessionHistoryResponse "OK"
// @Failure 400 {object} send.ErrorResponse "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /session/{sessionId}/history [get]
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

	isMember, err := h.sessionReader.UserIsMemberOfSession(r.Context(), sessionID, c.Subject)
	if err != nil {
		send.InternalServerError(w, "There has been an error fetching the session history")
		return
	}
	if !isMember {
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
