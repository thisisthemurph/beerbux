package url

import (
	"github.com/google/uuid"
	"net/http"
)

// GetUUIDFromPath gets a uuid.UUID value from the path value
func GetUUIDFromPath(r *http.Request, key string) (uuid.UUID, bool) {
	v, err := uuid.Parse(r.PathValue(key))
	if err != nil {
		return uuid.Nil, false
	}
	return v, true
}

func GetStringFromPath(r *http.Request, key string) (string, bool) {
	value := r.PathValue(key)
	if value == "" {
		return "", false
	}
	return value, true
}
