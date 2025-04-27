package middleware

import (
	"net/http"

	"example.com/goapi/pkg/httpx"
)

// Middleware
func AllowAccess(roles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the userID from the context
			// This is added by the Authentication middleware
			userData, ok := GetUserDetailsFromContext(r.Context())
			if !ok {
				httpx.Error(w, "user not authenticated", http.StatusUnauthorized)
				return
			}

			// Check if the user has at least one of the required roles
			hasRole := false
			for _, userRole := range userData.Roles {
				for _, requiredRole := range roles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}

				if hasRole {
					break
				}
			}

			// If user doesn't have required role, return Unauthorized
			if !hasRole {
				httpx.Error(w, "access denied", http.StatusForbidden)
				return
			}

			// User has the required role, proceed with the next handler
			next.ServeHTTP(w, r)
		})
	}
}
