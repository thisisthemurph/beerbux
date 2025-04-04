package user

import (
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/dto"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"net/http"
)

type GetCurrentUserHandler struct {
	userClient userpb.UserClient
}

func NewGetCurrentUserHandler(userClient userpb.UserClient) http.Handler {
	return &GetCurrentUserHandler{
		userClient: userClient,
	}
}

func (h *GetCurrentUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u, err := h.userClient.GetUser(r.Context(), &userpb.GetUserRequest{
		UserId: c.Subject,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	send.JSON(w, dto.UserResponse{
		ID:       u.UserId,
		Name:     u.Name,
		Username: u.Username,
	}, http.StatusOK)
}
