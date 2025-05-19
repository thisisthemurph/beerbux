package middleware

import (
	"log/slog"
	"net/http"
)

type RecoverMiddleware struct {
	logger *slog.Logger
}

func NewRecoverMiddleware(logger *slog.Logger) *RecoverMiddleware {
	return &RecoverMiddleware{
		logger: logger,
	}
}

// Recover handles recovering from a panic.
func (mw *RecoverMiddleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				mw.logger.Error("recovered from panic", "URL", r.URL.String(), "error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
