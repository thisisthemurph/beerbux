package url

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
)

// GetIDFromPath gets the value from the path value and sets a BadRequest status if the ID is not a UUID.
func GetIDFromPath(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	id, err := uuid.Parse(r.PathValue(key))
	if err != nil {
		send.Error(w, fmt.Sprintf("unable to parse id from path value: %s", key), http.StatusBadRequest)
		return "", false
	}

	return id.String(), true
}

func GetTextFromPath(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	value := r.PathValue(key)
	if value == "" {
		send.Error(w, fmt.Sprintf("unable to parse text from path value: %s", key), http.StatusBadRequest)
		return "", false
	}

	return value, value != ""
}
