package auth

import (
	"net/http"

	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/cookie"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers"
)

// RefreshHandler handles the refresh access token request.
// It expects a valid refresh token cookie and a valid access token cookie.
type RefreshHandler struct {
	authClient authpb.AuthClient
}

func NewRefreshHandler(authClient authpb.AuthClient) *RefreshHandler {
	return &RefreshHandler{
		authClient: authClient,
	}
}

func (h *RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshToken, ok := claims.GetRefreshToken(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := h.authClient.RefreshToken(r.Context(), &authpb.RefreshTokenRequest{
		UserId:       c.Subject,
		RefreshToken: refreshToken,
	})

	if err != nil {
		handlers.WriteValidationError(w, err)
		return
	}

	cookie.SetAccessTokenCookie(w, user.AccessToken)
	cookie.SetRefreshTokenCookie(w, user.RefreshToken)

	w.WriteHeader(http.StatusNoContent)
}
