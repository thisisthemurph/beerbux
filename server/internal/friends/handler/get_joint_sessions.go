package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/friends/db"
	"beerbux/internal/friends/query"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
	"log/slog"
	"net/http"
	"time"
)

type GetJointSessionsHandler struct {
	membersAreFriendsQuery *query.MembersAreFriendsQuery
	getJointSessionsQuery  *query.GetJointSessionsQuery
	logger                 *slog.Logger
}

func NewGetJointSessionsHandler(
	membersAreFriendsQuery *query.MembersAreFriendsQuery,
	getJointSessionsQuery *query.GetJointSessionsQuery,
	logger *slog.Logger,
) *GetJointSessionsHandler {
	return &GetJointSessionsHandler{
		membersAreFriendsQuery: membersAreFriendsQuery,
		getJointSessionsQuery:  getJointSessionsQuery,
		logger:                 logger,
	}
}

type SessionResultItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatorID uuid.UUID `json:"creatorID"`
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

	ss, err := h.getJointSessionsQuery.Execute(r.Context(), c.Subject, friendID)
	if err != nil {
		h.logger.Error("failed to fetch sessions", "error", err)
		send.InternalServerError(w, "There has been an issue fetching your shared sessions")
		return
	}

	results := fn.Map(ss, func(s db.Session) SessionResultItem {
		return SessionResultItem{
			ID:        s.ID,
			Name:      s.Name,
			IsActive:  s.IsActive,
			CreatorID: s.CreatorID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		}
	})

	send.JSON(w, results, http.StatusOK)
}
