package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/session/command"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type LeaveSessionHandler struct {
	removeSessionMemberCommand *command.RemoveSessionMemberCommand
	logger                     *slog.Logger
}

func NewLeaveSessionHandler(removeSessionMemberCommand *command.RemoveSessionMemberCommand, logger *slog.Logger) *LeaveSessionHandler {
	return &LeaveSessionHandler{
		removeSessionMemberCommand: removeSessionMemberCommand,
		logger:                     logger,
	}
}

// LeaveSessionHandler godoc
// @Summary Leave Session
// @Tags session
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 "OK"
// @Failure 400 {object} send.ErrorResponse "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /session/{sessionId}/leave [post]
func (h *LeaveSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, ok := url.Path.GetUUID(r, "sessionId")
	if !ok {
		send.BadRequest(w, "Session ID is required")
		return
	}

	err := h.removeSessionMemberCommand.Execute(r.Context(), sessionID, c.Subject, c.Subject)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "Session not found")
		} else if errors.Is(err, sessionErr.ErrSessionMemberNotFound) {
			send.NotFound(w, "Session member not found")
		} else if errors.Is(err, sessionErr.ErrSessionMustHaveAtLeastOneMember) {
			send.BadRequest(w, "You cannot leave the session if you are the only member")
		} else if errors.Is(err, sessionErr.ErrSessionMustHaveAtLeastOneAdmin) {
			send.BadRequest(w, "You cannot leave the session if you are the only admin member")
		} else {
			h.logger.Error("failed to remove member from session", "error", err)
			send.InternalServerError(w, "There was an issue leaving the session")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

type LeaveSessionHandlerURLParams struct {
	SessionID uuid.UUID
}

func (h *LeaveSessionHandler) getURLParams(r *http.Request) (LeaveSessionHandlerURLParams, bool) {
	sessionID, ok := url.Path.GetUUID(r, "sessionId")
	if !ok {
		return LeaveSessionHandlerURLParams{}, false
	}

	return LeaveSessionHandlerURLParams{
		SessionID: sessionID,
	}, true
}
