package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/cookie"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type LoginHandler struct {
	generateTokensCommand  *command.GenerateTokensCommand
	comparePasswordCommand *command.ComparePasswordCommand
	logger                 *slog.Logger
}

func NewLoginHandler(
	loginCommand *command.GenerateTokensCommand,
	comparePasswordCommand *command.ComparePasswordCommand,
	logger *slog.Logger,
) *LoginHandler {
	return &LoginHandler{
		generateTokensCommand:  loginCommand,
		comparePasswordCommand: comparePasswordCommand,
		logger:                 logger,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		send.ValidationError(w, err)
		return
	}

	if err := h.comparePasswordCommand.Execute(r.Context(), req.Username, req.Password); err != nil {
		if errors.Is(err, command.ErrPasswordMismatch) || errors.Is(err, command.ErrUserNotFound) {
			send.Unauthorized(w, "Invalid username or password")
		} else {
			h.logger.Error("failed when comparing login password", "error", err)
			send.InternalServerError(w, "There was an issue logging in")
		}
		return
	}

	tokens, err := h.generateTokensCommand.Execute(r.Context(), req.Username)
	if err != nil {
		h.handleLoginError(w, err)
		return
	}

	cookie.SetAccessTokenCookie(w, tokens.AccessToken)
	cookie.SetRefreshTokenCookie(w, tokens.RefreshToken)

	send.JSON(w, LoginResponse{
		ID:       tokens.User.ID,
		Name:     tokens.User.Name,
		Username: tokens.User.Username,
	}, http.StatusOK)
}

func (h *LoginHandler) handleLoginError(w http.ResponseWriter, err error) {
	if errors.Is(err, command.ErrUserNotFound) {
		send.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	h.logger.Error("error signing in", "error", err)
	send.Error(w, "There was an issue signing you in", http.StatusInternalServerError)
}

func (r LoginRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Username, oz.Required),
		oz.Field(&r.Password, oz.Required),
	)
}
