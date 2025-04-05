package user

import (
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type GetBalanceHandler struct {
	userClient userpb.UserClient
}

func NewGetBalanceHandler(userClient userpb.UserClient) http.Handler {
	return &GetBalanceHandler{
		userClient: userClient,
	}
}

type BalanceResponse struct {
	Credit float64 `json:"credit"`
	Debit  float64 `json:"debit"`
	Net    float64 `json:"net"`
}

func (g GetBalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Determine who should be able to see the balance of another user
	// Should these user's be friends before they can see each other's balance?
	// For now allow any user to see any other user's balance, but the user must be authed.
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance, err := g.userClient.GetUserBalance(r.Context(), &userpb.GetUserRequest{
		UserId: userID.String(),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	send.JSON(w, BalanceResponse{
		Credit: balance.Credit,
		Debit:  balance.Debit,
		Net:    balance.Net,
	}, http.StatusOK)
}
