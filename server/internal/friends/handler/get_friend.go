package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/useraccess"
	"beerbux/internal/friends/query"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"log/slog"
	"net/http"
)

type GetFriendHandler struct {
	userReader             useraccess.UserReader
	membersAreFriendsQuery *query.MembersAreFriendsQuery
	logger                 *slog.Logger
}

func NewGetFriendHandler(userReader useraccess.UserReader, friendsQuery *query.MembersAreFriendsQuery, logger *slog.Logger) *GetFriendHandler {
	return &GetFriendHandler{
		userReader:             userReader,
		membersAreFriendsQuery: friendsQuery,
		logger:                 logger,
	}
}

func (h *GetFriendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	friendID, ok := url.Path.GetUUID(r, "friendId")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	areFriends, err := h.membersAreFriendsQuery.Execute(r.Context(), c.Subject, friendID)
	if err != nil {
		h.logger.Error("error determining if users are friends", "error", err, "userID", c.Subject, "friendID", friendID)
		send.InternalServerError(w, "Failed to determine if this user is your friend.")
		return
	}
	if !areFriends {
		send.Unauthorized(w, "You are not friends with this user")
		return
	}

	friend, err := h.userReader.GetUserByID(r.Context(), friendID)
	if err != nil {
		if errors.Is(err, useraccess.ErrUserNotFound) {
			send.NotFound(w, "Friend not found")
		} else {
			h.logger.Error("error fetching user", "error", err, "friendID", friendID)
			send.InternalServerError(w, "Failed to fetch friend")
		}
		return
	}

	send.JSON(w, friend, http.StatusOK)
}
