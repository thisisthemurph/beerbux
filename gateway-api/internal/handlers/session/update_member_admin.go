package session

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
)

type UpdateSessionMemberAdmin struct {
	logger        *slog.Logger
	sessionClient sessionpb.SessionClient
}

func NewUpdateSessionMemberAdmin(logger *slog.Logger, sessionClient sessionpb.SessionClient) *UpdateSessionMemberAdmin {
	return &UpdateSessionMemberAdmin{
		logger:        logger,
		sessionClient: sessionClient,
	}
}

type UpdateSessionMemberAdminRequest struct {
	NewAdminState bool `json:"newAdminState"`
}

func (h *UpdateSessionMemberAdmin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("sessionId"))
	if err != nil {
		h.logger.Error("error parsing sessionId", "error", err)
		send.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	memberID, err := uuid.Parse(r.PathValue("memberId"))
	if err != nil {
		h.logger.Error("error parsing memberId", "error", err)
		send.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	var req UpdateSessionMemberAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("error decoding body", "error", err)
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	if memberID.String() == c.Subject {
		send.Error(w, "You cannot update your own admin status", http.StatusBadRequest)
		return
	}

	ssn, err := h.sessionClient.GetSession(r.Context(), &sessionpb.GetSessionRequest{
		SessionId: sessionID.String(),
	})
	if err != nil {
		h.logger.Error("error getting session", "error", err)
		send.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	var currentMember *sessionpb.SessionMember
	var updateMember *sessionpb.SessionMember
	for _, m := range ssn.Members {
		if m.UserId == c.Subject {
			currentMember = m
		}
		if m.UserId == memberID.String() {
			updateMember = m
		}
		if currentMember != nil && updateMember != nil {
			break
		}
	}

	if currentMember == nil || !currentMember.IsAdmin {
		send.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if updateMember == nil {
		send.Error(w, "Member to update not found", http.StatusNotFound)
		return
	}

	if req.NewAdminState == updateMember.IsAdmin {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err = h.sessionClient.UpdateSessionMemberAdminState(r.Context(), &sessionpb.UpdateSessionMemberAdminStateRequest{
		SessionId: sessionID.String(),
		UserId:    memberID.String(),
		IsAdmin:   req.NewAdminState,
	})
	if err != nil {
		h.logger.Error("error updating session member", "error", err)

		sc, ok := status.FromError(err)
		if ok {
			switch sc.Code() {
			case codes.FailedPrecondition:
				send.Error(w, "A session must have at least one admin member.", http.StatusBadRequest)
			case codes.Internal:
				send.Error(w, "Unexpected error updating the admin status of the member", http.StatusInternalServerError)
			}
			return
		}
		send.Error(w, "There was an unexpected error updating the admin state of the session member", http.StatusInternalServerError)
		return
	}
}
