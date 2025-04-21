package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func LogContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a logger with request-specific fields
		logger := log.With().
			Str("request_id", middleware.GetReqID(r.Context())).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Logger()

		// Inject logger into context
		ctx := logger.WithContext(r.Context())

		// Chi's response wrapper to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r.WithContext(ctx))

		// Log request completion
		logger.Info().
			Int("status", ww.Status()).
			Dur("duration_ms", time.Since(start)).
			Msg("request completed")
	})
}
