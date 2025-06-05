package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
)

type ListCurrentUserSessionsHandler struct {
	listSessionsByUserIDQuery *query.ListSessionsByUserIDQuery
	logger                    *slog.Logger
}

func NewListCurrentUserSessionsHandler(listSessionsByUserIDQuery *query.ListSessionsByUserIDQuery, logger *slog.Logger) *ListCurrentUserSessionsHandler {
	return &ListCurrentUserSessionsHandler{
		listSessionsByUserIDQuery: listSessionsByUserIDQuery,
		logger:                    logger,
	}
}

type CurrentUserSessionResponse struct {
	ID       uuid.UUID                  `json:"id"`
	Name     string                     `json:"name"`
	Total    float64                    `json:"total"`
	IsActive bool                       `json:"isActive"`
	Members  []CurrentUserSessionMember `json:"members"`
}

type CurrentUserSessionMember struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
}

// LeaveSessionHandler godoc
// @Summary List Current User Sessions
// @Tags session
// @Accept json
// @Produce json
// @Success 200 {array} CurrentUserSessionResponse "OK"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /user/session [get]
func (h *ListCurrentUserSessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	limit := 0
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if i, err := strconv.Atoi(ps); err == nil {
			limit = i
		}
	}

	ss, err := h.listSessionsByUserIDQuery.Execute(r.Context(), c.Subject, int32(limit))
	if err != nil {
		h.logger.Error("error listing sessions for user", "user", c.Subject, "error", err)
		send.InternalServerError(w, "There has been an issue listing your sessions")
		return
	}

	sessions := make([]CurrentUserSessionResponse, 0, len(ss))
	for _, s := range ss {
		mm := make([]CurrentUserSessionMember, 0, len(s.Members))
		for _, m := range s.Members {
			mm = append(mm, CurrentUserSessionMember{
				ID:       m.ID,
				Name:     m.Name,
				Username: m.Username,
			})
		}

		sessions = append(sessions, CurrentUserSessionResponse{
			ID:       s.ID,
			Name:     s.Name,
			Total:    s.Total,
			IsActive: s.IsActive,
			Members:  mm,
		})
	}

	send.JSON(w, sessions, http.StatusOK)
}
