package auth

import (
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/cookie"
	"net/http"
	"time"
)

type LogoutHandler struct {
	authClient authpb.AuthClient
}

func NewLogoutHandler(authClient authpb.AuthClient) *LogoutHandler {
	return &LogoutHandler{
		authClient: authClient,
	}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Invalidate the refresh token if it exists.
	refreshToken, hasRefreshToken := claims.GetRefreshToken(r)
	if hasRefreshToken {
		_, _ = h.authClient.InvalidateRefreshToken(r.Context(), &authpb.InvalidateRefreshTokenRequest{
			UserId:       c.Subject,
			RefreshToken: refreshToken,
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookie.AccessTokenKey,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     cookie.RefreshTokenKey,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}
