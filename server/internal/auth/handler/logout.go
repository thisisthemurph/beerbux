package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/cookie"
	"beerbux/internal/common/claims"
	"log/slog"
	"net/http"
	"time"
)

type LogoutHandler struct {
	invalidateRefreshTokenCommand *command.InvalidateRefreshTokenCommand
	logger                        *slog.Logger
}

func NewLogoutHandler(invalidateRefreshTokenCommand *command.InvalidateRefreshTokenCommand, logger *slog.Logger) *LogoutHandler {
	return &LogoutHandler{
		invalidateRefreshTokenCommand: invalidateRefreshTokenCommand,
		logger:                        logger,
	}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	refreshToken, ok := claims.GetRefreshToken(r)
	if ok {
		if err := h.invalidateRefreshTokenCommand.Execute(r.Context(), c.Subject, refreshToken); err != nil {
			h.logger.Error("failed to invalidate refresh token", "user", c.Subject, "error", err)
		}
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
