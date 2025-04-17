package middleware

import (
	"log"
	"net/http"
	"time"
)

// All middleware returns this type
type Middleware func(next http.Handler) http.Handler

// NOTE:
// ✅ Use simple direct style when the middleware has no configuration/state.
//     eg: func Recoverer(next http.Handler) http.Handler {
// 		       fn:=(some function)
//             return http.HandlerFunc(fn)
// 		   }
//     - for middleware that don't need configurations
//
// ✅ Use the func() func() (middleware factory) style is used when you want to parameterize or inject config into your middleware.
//     eg: the `LogTime` middleware

func LogTime() Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Printf("➡️  %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
			log.Printf("⏱️  Completed in %v", time.Since(start))
		}

		return http.HandlerFunc(fn)
	}
}
