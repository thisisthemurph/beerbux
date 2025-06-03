package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/friends/query"
	"beerbux/pkg/send"
	"log/slog"
	"net/http"
)

type GetFriendsHandler struct {
	getFriendsQuery *query.GetFriendsQuery
	logger          *slog.Logger
}

func NewGetFriendsHandler(getFriendsQuery *query.GetFriendsQuery, logger *slog.Logger) *GetFriendsHandler {
	return &GetFriendsHandler{
		getFriendsQuery: getFriendsQuery,
		logger:          logger,
	}
}

// GetFriendsHandler godoc
// @Summary Get Friends
// @Tags friends
// @Accept json
// @Produce json
// @Success 200 {array} query.Friend
// @Failure 401 "Unauthorized"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /friends [get]
func (h *GetFriendsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	friends, err := h.getFriendsQuery.Execute(r.Context(), c.Subject)
	if err != nil {
		send.InternalServerError(w, "There has been an issue fetching your list of friends")
		return
	}

	send.JSON(w, friends, http.StatusOK)
}
