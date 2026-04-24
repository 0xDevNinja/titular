package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Recovery catches panics and returns a generic 500 without leaking internals.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().
					Str("request_id", GetRequestID(r.Context())).
					Interface("panic", rec).
					Msg("recovered from panic")
				http.Error(w, `{"error":{"code":"internal_server_error","message":"an unexpected error occurred"}}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
