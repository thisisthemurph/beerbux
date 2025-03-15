package auth

import (
	"encoding/json"
	"net/http"

	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers"
)

type LoginHandler struct {
	authClient authpb.AuthClient
}

func NewLoginHandler(authClient authpb.AuthClient) *LoginHandler {
	return &LoginHandler{
		authClient: authClient,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		handlers.WriteValidationError(w, err)
		return
	}

	user, err := h.authClient.Login(r.Context(), &authpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		handlers.WriteValidationError(w, err)
		return
	}

	setAccessTokenCookie(w, user.AccessToken)
	setRefreshTokenCookie(w, user.RefreshToken)

	w.WriteHeader(http.StatusOK)
}

func (r LoginRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Username, oz.Required),
		oz.Field(&r.Password, oz.Required),
	)
}
