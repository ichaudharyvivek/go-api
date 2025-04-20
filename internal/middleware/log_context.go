package middleware

import (
	"net/http"
	"time"

	"example.com/goapi/internal/utils/logger"
	"github.com/go-chi/chi/v5/middleware"

	"go.uber.org/zap"
)

func LogContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get request ID from Chi's middleware
		requestID := r.Context().Value(middleware.RequestIDKey).(string)

		// Create context with logger fields
		ctx := logger.NewContext(r.Context(),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("request_id", requestID),
		)

		// Continue with enriched context
		next.ServeHTTP(w, r.WithContext(ctx))

		// Log completion
		logger.FromContext(ctx).Info("Request completed",
			zap.Duration("duration", time.Since(start)),
		)
	}

	return http.HandlerFunc(fn)
}
