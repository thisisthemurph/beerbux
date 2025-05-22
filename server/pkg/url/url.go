package url

import (
	"github.com/google/uuid"
	"net/http"
)

type path struct{}

var Path path = path{}

// GetUUID returns a uuid.UUID value from the path value for the given key and a boolean
// indicating if the uuid.UUID was present and valid.
func (p path) GetUUID(r *http.Request, key string) (uuid.UUID, bool) {
	v, err := uuid.Parse(r.PathValue(key))
	if err != nil {
		return uuid.Nil, false
	}
	return v, true
}

// GetString returns a string value from the path and a boolean indicating if that string was present.
func (p path) GetString(r *http.Request, key string) (string, bool) {
	value := r.PathValue(key)
	if value == "" {
		return "", false
	}
	return value, true
}

type query struct{}

var Query query = query{}

// GetString returns a string value from the URL query params and a boolean indicating if that string was present.
func (q query) GetString(r *http.Request, key string) (string, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return "", false
	}
	return value, true
}
