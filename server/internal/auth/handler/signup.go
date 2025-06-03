package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"fmt"
	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"log/slog"
	"net/http"
)

var ErrPasswordsDoNotMatch = errors.New("passwords do not match")

type SignupHandler struct {
	signupCommand *command.SignupCommand
	logger        *slog.Logger
}

func NewSignupHandler(signupCommand *command.SignupCommand, logger *slog.Logger) *SignupHandler {
	return &SignupHandler{
		signupCommand: signupCommand,
		logger:        logger,
	}
}

type SignupRequest struct {
	Name                 string `json:"name"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	VerificationPassword string `json:"verificationPassword"`
}

// SignupHandler godoc
// @Summary Signup
// @Tags auth
// @Accept json
// @Param login body SignupRequest true "Login credentials"
// @Success 201
// @Failure 400 {object} send.ErrorResponse "Bad Request"
// @Failure 400 {object} send.ValidationErrorModel "Bad Request (validation error)"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("bad request for signup", "error", err)
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		send.ValidationError(w, err)
		return
	}

	_, err := h.signupCommand.Execute(r.Context(), req.Name, req.Username, req.Email, req.Password, req.VerificationPassword)
	if err != nil {
		h.handleSignupError(w, req, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SignupHandler) handleSignupError(w http.ResponseWriter, req SignupRequest, err error) {
	if errors.Is(err, command.ErrPasswordMismatch) {
		send.Error(w, "The provided passwords do not match", http.StatusBadRequest)
		return
	}
	if errors.Is(err, command.ErrUsernameTaken) {
		send.Error(w, fmt.Sprintf("Username %s already taken", req.Username), http.StatusBadRequest)
		return
	}

	h.logger.Error("signup command error", "error", err)
	send.Error(w, "There has been an issue creating your account, please try again", http.StatusBadRequest)
}

func (r SignupRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Name,
			oz.Required.Error("Name is required"),
			oz.Length(2, 50).Error("Name must be between 2 and 50 characters")),
		oz.Field(&r.Username,
			oz.Required.Error("Username is required"),
			oz.Length(3, 25).Error("Username must be between 3 and 25 characters")),
		oz.Field(&r.Email,
			oz.Required.Error("Email address is required"),
			is.Email.Error("The provided email is not a valid email address")),
		oz.Field(&r.Password,
			oz.Required.Error("Password is required"),
			oz.Length(8, 0).Error("Password must be at least 8 characters")),
		oz.Field(&r.VerificationPassword,
			oz.Required.Error("Verification password is required")),
		oz.Field(&r.Password,
			oz.By(func(value interface{}) error {
				if r.Password != r.VerificationPassword {
					return ErrPasswordsDoNotMatch
				}
				return nil
			}),
		),
	)
}
