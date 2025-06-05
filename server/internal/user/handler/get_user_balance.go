package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/useraccess"
	"beerbux/pkg/send"
	"errors"
	"log/slog"
	"net/http"
)

type GetUserBalanceHandler struct {
	userReader useraccess.UserReader
	logger     *slog.Logger
}

func NewGetCurrentUserBalanceHandler(userReader useraccess.UserReader, logger *slog.Logger) *GetUserBalanceHandler {
	return &GetUserBalanceHandler{
		userReader: userReader,
		logger:     logger,
	}
}

type BalanceResponse struct {
	Credit float64 `json:"credit"`
	Debit  float64 `json:"debit"`
	Net    float64 `json:"net"`
}

// GetCurrentUserHandler godoc
// @Summary Get Current User Balance
// @Tags user
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} BalanceResponse "OK"
// @Failure 401 "Unauthorized"
// @Failure 404 {object} send.ErrorResponse "Not Found"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /user/{userId}/balance [get]
func (h *GetUserBalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	send.JSON(w, BalanceResponse{
		Credit: user.Account.Credit,
		Debit:  user.Account.Debit,
		Net:    user.Account.Credit - user.Account.Debit,
	}, http.StatusOK)
}
