package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/useraccess"
	"beerbux/pkg/send"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type GetCurrentUserHandler struct {
	userReader useraccess.UserReader
	logger     *slog.Logger
}

func NewGetCurrentUserHandler(userReader useraccess.UserReader, logger *slog.Logger) *GetCurrentUserHandler {
	return &GetCurrentUserHandler{
		userReader: userReader,
		logger:     logger,
	}
}

type GetCurrentUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
}

// GetCurrentUserHandler godoc
// @Summary Get Current User
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} GetCurrentUserResponse "OK"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} send.ErrorResponse "Not Found"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /user [get]
func (h *GetCurrentUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := h.userReader.GetUserByID(r.Context(), c.Subject)
	if err != nil {
		if errors.Is(useraccess.ErrUserNotFound, err) {
			send.NotFound(w, "Your user account could not be found")
			return
		}
		h.logger.Error("error getting user", "user", c.Subject, "error", err)
		send.InternalServerError(w, "There was an issue finding your user account")
		return
	}

	send.JSON(w, GetCurrentUserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}, http.StatusOK)
}
