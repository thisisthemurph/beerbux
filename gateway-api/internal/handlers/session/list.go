package session

import (
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/pagination"
	"net/http"

	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
)

type ListSessionsForUserHandler struct {
	sessionClient sessionpb.SessionClient
}

func NewListSessionsForUserHandler(sessionClient sessionpb.SessionClient) *ListSessionsForUserHandler {
	return &ListSessionsForUserHandler{
		sessionClient: sessionClient,
	}
}

type ListSessionsForUserResponse struct {
	ID       string                      `json:"id"`
	Name     string                      `json:"name"`
	IsActive bool                        `json:"isActive"`
	Members  []ListSessionsForUserMember `json:"members"`
}

type ListSessionsForUserMember struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (h *ListSessionsForUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	page := pagination.FromRequest(r)
	ssns, err := h.sessionClient.ListSessionsForUser(r.Context(), &sessionpb.ListSessionsForUserRequest{
		UserId:    c.Subject,
		PageSize:  page.PageSize,
		PageToken: page.PageToken,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessions := make([]ListSessionsForUserResponse, 0, len(ssns.Sessions))
	for _, ssn := range ssns.Sessions {
		mm := make([]ListSessionsForUserMember, 0, len(ssn.Members))
		for _, m := range ssn.Members {
			mm = append(mm, ListSessionsForUserMember{
				ID:       m.UserId,
				Name:     m.Name,
				Username: m.Username,
			})
		}

		sessions = append(sessions, ListSessionsForUserResponse{
			ID:       ssn.SessionId,
			Name:     ssn.Name,
			IsActive: ssn.IsActive,
			Members:  mm,
		})
	}

	send.JSON(w, sessions, http.StatusOK)
}
