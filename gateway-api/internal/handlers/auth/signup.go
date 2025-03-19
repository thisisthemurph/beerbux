package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SignupHandler struct {
	authClient authpb.AuthClient
}

func NewSignupHandler(authClient authpb.AuthClient) http.Handler {
	return &SignupHandler{
		authClient: authClient,
	}
}

var ErrPasswordsDoNotMatch = errors.New("passwords do not match")

type SignupRequest struct {
	Name                 string `json:"name"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	VerificationPassword string `json:"verificationPassword"`
}

func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		handlers.WriteValidationError(w, err)
		return
	}

	_, err := h.authClient.Signup(r.Context(), &authpb.SignupRequest{
		Name:                 req.Name,
		Username:             req.Username,
		Password:             req.Password,
		VerificationPassword: req.VerificationPassword,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteError(w, "There has been an error logging you up", http.StatusUnauthorized)
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			// Username is already taken or passwords do not match.
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteError(w, st.Message(), http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteError(w, "There has been an error logging you up", http.StatusUnauthorized)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (r SignupRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Name, oz.Required.Error("Name is required"), oz.Length(2, 50).Error("Name must be between 2 and 50 characters")),
		oz.Field(&r.Username, oz.Required.Error("Username is required"), oz.Length(3, 25).Error("Username must be between 3 and 25 characters")),
		oz.Field(&r.Password, oz.Required.Error("Password is required"), oz.Length(8, 0).Error("Password must be at least 8 characters")),
		oz.Field(&r.VerificationPassword, oz.Required.Error("Verification password is required")),
		oz.Field(&r.Password, oz.By(func(value interface{}) error {
			if r.Password != r.VerificationPassword {
				return ErrPasswordsDoNotMatch
			}
			return nil
		}),
		),
	)
}
