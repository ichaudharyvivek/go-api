package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/goapi/internal/common/errors"
	"example.com/goapi/internal/domain/auth"
	m "example.com/goapi/internal/middleware"
	_v "example.com/goapi/internal/utils/validator"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	service   auth.Service
	validator *validator.Validate
}

func NewAuthHandler(s auth.Service, v *validator.Validate) *AuthHandler {
	return &AuthHandler{service: s, validator: v}
}

func (h *AuthHandler) RegisterAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(m.Authenticate("secret"))
			r.Post("/logout", h.Logout)
			r.Post("/refresh", h.RefreshTokens)
		})
	})
}

// RegisterUser handles user registration
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	payload := &auth.RegisterUserPayload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		httpx.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		respBody := _v.ToErrResponse(err)
		httpx.Errors(w, respBody, http.StatusBadRequest)
		return
	}

	user, err := h.service.Create(r.Context(), payload)
	if err != nil {
		if apiErr, ok := err.(*errors.ApiError); ok {
			logger.Error().Err(err).Str("path", r.URL.Path).Str("method", r.Method).Msg("Failed to create user")
			httpx.Error(w, apiErr.Message, http.StatusInternalServerError)
			return
		}

		logger.Error().Err(err).Msg("Unexpected error during user registration")
		httpx.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Assuming user.ToDto() exists and returns a sanitized user DTO
	httpx.Created(w, user.ToDto())
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	var credentials struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		httpx.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(credentials); err != nil {
		respBody := _v.ToErrResponse(err)
		httpx.Errors(w, respBody, http.StatusBadRequest)
		return
	}

	tokens, err := h.service.Login(r.Context(), credentials.Email, credentials.Password)
	if err != nil {
		if apiErr, ok := err.(*errors.ApiError); ok {
			logger.Error().Err(err).Msg("Login failed")
			httpx.Error(w, apiErr.Message, http.StatusInternalServerError)
			return
		}

		logger.Error().Err(err).Msg("Unexpected error during login")
		httpx.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set refresh token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Enable in production (HTTPS only)
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(1 * 24 * time.Hour), // 1 day
		// MaxAge:   int(24 * time.Hour * 7), // This is redundant but can be used together with expires because some browsers still check this attribute
	})

	httpx.Ok(w, map[string]any{
		"access_token": tokens.AccessToken,
		"expires_in":   "15 mins",
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	// Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httpx.Error(w, "Missing refresh token", http.StatusBadRequest)
		return
	}
	// If cookie exists, log its value
	logger.Debug().Str("refresh_token", cookie.Value).Msg("Retrieved refresh token")

	if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
		if apiErr, ok := err.(*errors.ApiError); ok {
			logger.Error().Err(err).Msg("Logout failed")
			httpx.Error(w, apiErr.Message, http.StatusInternalServerError)
			return
		}

		logger.Error().Err(err).Msg("Unexpected error during logout")
		httpx.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Clear the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Immediately expire the cookie
	})
	httpx.Ok(w, map[string]string{"message": "Successfully logged out"})
}

// RefreshTokens handles token refresh
func (h *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	// Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httpx.Error(w, "Missing refresh token", http.StatusBadRequest)
		return
	}

	tokens, err := h.service.RefreshTokens(r.Context(), cookie.Value)
	if err != nil {
		if apiErr, ok := err.(*errors.ApiError); ok {
			logger.Error().Err(err).Msg("Token refresh failed")
			httpx.Error(w, apiErr.Message, http.StatusInternalServerError)
			return
		}

		logger.Error().Err(err).Msg("Unexpected error during token refresh")
		httpx.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set new refresh token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(24 * time.Hour * 7), // 7 days
	})

	httpx.Ok(w, map[string]interface{}{
		"access_token": tokens.AccessToken,
		"expires_in":   "15 mins",
	})
}
