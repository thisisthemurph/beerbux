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
	loginCommand *command.LoginCommand
	logger       *slog.Logger
}

func NewLoginHandler(loginCommand *command.LoginCommand, logger *slog.Logger) *LoginHandler {
	return &LoginHandler{
		loginCommand: loginCommand,
		logger:       logger,
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

	resp, err := h.loginCommand.Execute(r.Context(), req.Username, req.Password)
	if err != nil {
		h.handleLoginError(w, err)
		return
	}

	cookie.SetAccessTokenCookie(w, resp.AccessToken)
	cookie.SetRefreshTokenCookie(w, resp.RefreshToken)

	send.JSON(w, LoginResponse{
		ID:       resp.User.ID,
		Name:     resp.User.Name,
		Username: resp.User.Username,
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
