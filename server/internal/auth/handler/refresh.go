package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/cookie"
	"beerbux/internal/common/claims"
	"errors"
	"log/slog"
	"net/http"
)

type RefreshHandler struct {
	refreshTokenCommand *command.RefreshTokenCommand
	logger              *slog.Logger
}

func NewRefreshHandler(refreshTokenCommand *command.RefreshTokenCommand, logger *slog.Logger) *RefreshHandler {
	return &RefreshHandler{
		refreshTokenCommand: refreshTokenCommand,
		logger:              logger,
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

	resp, err := h.refreshTokenCommand.Execute(r.Context(), c.Subject, refreshToken)
	if err != nil {
		if errors.Is(err, command.ErrRefreshTokenNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie.SetAccessTokenCookie(w, resp.AccessToken)
	cookie.SetRefreshTokenCookie(w, resp.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}
