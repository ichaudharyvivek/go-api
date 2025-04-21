package middleware

import (
	"context"
	"net/http"
	"strings"

	"example.com/goapi/pkg/httpx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authKey struct{}

var userIDKey = authKey{}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				httpx.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				// Verify signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				httpx.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				httpx.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok {
				httpx.Error(w, "token missing subject", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(sub)
			if err != nil {
				httpx.Error(w, "invalid user in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}

	return ""
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}
