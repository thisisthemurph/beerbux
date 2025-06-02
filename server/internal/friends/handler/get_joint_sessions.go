package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/sessionaccess"
	"beerbux/internal/friends/query"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type GetJointSessionsHandler struct {
	membersAreFriendsQuery *query.MembersAreFriendsQuery
	getJointSessionsQuery  *query.GetJointSessionIDsQuery
	sessionReader          sessionaccess.SessionReader
	logger                 *slog.Logger
}

func NewGetJointSessionsHandler(
	membersAreFriendsQuery *query.MembersAreFriendsQuery,
	getJointSessionsQuery *query.GetJointSessionIDsQuery,
	sessionReader sessionaccess.SessionReader,
	logger *slog.Logger,
) *GetJointSessionsHandler {
	return &GetJointSessionsHandler{
		membersAreFriendsQuery: membersAreFriendsQuery,
		getJointSessionsQuery:  getJointSessionsQuery,
		sessionReader:          sessionReader,
		logger:                 logger,
	}
}

type SessionResultItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatorID uuid.UUID `json:"creatorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h *GetJointSessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	friendID, found := url.Path.GetUUID(r, "friendId")
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	areFriends, err := h.membersAreFriendsQuery.Execute(r.Context(), c.Subject, friendID)
	if err != nil {
		h.logger.Error("failed to fetch shared sessions", "error", err)
		send.InternalServerError(w, "There has been an issue fetching your shared sessions")
		return
	} else if !areFriends {
		send.Unauthorized(w, "You are not friends with this member")
		return
	}

	jointSessionIDs, err := h.getJointSessionsQuery.Execute(r.Context(), c.Subject, friendID)
	if err != nil {
		h.logger.Error("failed to fetch sessions", "error", err)
		send.InternalServerError(w, "There has been an issue fetching your shared sessions")
		return
	}

	ss := make([]*sessionaccess.Session, 0, len(jointSessionIDs))
	for _, sessionID := range jointSessionIDs {
		s, err := h.sessionReader.GetSessionDetails(r.Context(), sessionID)
		if err != nil {
			h.logger.Error("failed to fetch session", "error", err)
			continue
		}
		ss = append(ss, s)
	}

	send.JSON(w, ss, http.StatusOK)
}
