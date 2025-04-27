package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"example.com/goapi/pkg/httpx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type userDataKey struct{}
type UserDataContext struct {
	UserID uuid.UUID
	Roles  []string
}

var userDataContext = userDataKey{}

// NOTE: The secret here is the same key we used to create accessToken.
// Use the same key or the middleware will fail with invalid signature exception.
// Authenticate is a middleware that checks if the request is authenticated via JWT.
func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				httpx.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			// Parse and validate the token
			token, err := parseAndValidateToken(tokenString, secret)
			if err != nil {
				httpx.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Extract claims from the token
			userID, roles, err := extractUserDataFromClaims(token)
			if err != nil {
				httpx.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Add the user data to the context
			userData := UserDataContext{
				UserID: userID,
				Roles:  roles,
			}

			ctx := context.WithValue(r.Context(), userDataContext, userData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Extracts the JWT token from the Authorization header.
func extractTokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}

	return ""
}

// Parses the JWT token and checks for expiration and validity.
func parseAndValidateToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}

		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

// Extracts userID and roles from JWT claims.
func extractUserDataFromClaims(token *jwt.Token) (uuid.UUID, []string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, nil, jwt.ErrTokenInvalidClaims
	}

	// Extract userID
	userIDStr, ok := claims["userID"].(string)
	if !ok {
		return uuid.UUID{}, nil, errors.New("token missing userID")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, nil, errors.New("invalid userID in token")
	}

	rolesInterface, ok := claims["roles"].([]interface{})
	if !ok {
		return uuid.UUID{}, nil, errors.New("token missing roles or invalid format")
	}

	roles := make([]string, 0, len(rolesInterface))
	for _, role := range rolesInterface {
		if roleStr, ok := role.(string); ok {
			roles = append(roles, roleStr)
		} else {
			return uuid.UUID{}, nil, errors.New("role in token is not a string")
		}
	}

	return userID, roles, nil
}

// Extract User Details from Context. Includes UserID and Roles
func GetUserDetailsFromContext(ctx context.Context) (UserDataContext, bool) {
	userData, ok := ctx.Value(userDataContext).(UserDataContext)
	return userData, ok
}

// Prev function where we used just authKey struct{}
// func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
// 	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
// 	return userID, ok
// }
