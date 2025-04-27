package auth

import (
	"context"
	"time"

	"example.com/goapi/internal/common/errors"
	"example.com/goapi/internal/domain/user"
	m "example.com/goapi/internal/middleware"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, payload *RegisterUserPayload) (*user.User, error)
	Login(ctx context.Context, email, password string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error)
}

type service struct {
	repo         Repository
	accessToken  TokenConfig
	refreshToken TokenConfig
}

func NewService(r Repository, accessSecret, refreshSecret string, accessExp, refreshExp time.Duration) Service {
	return &service{
		repo: r,
		accessToken: TokenConfig{
			Secret:     accessSecret,
			Expiration: accessExp,
		},
		refreshToken: TokenConfig{
			Secret:     refreshSecret,
			Expiration: refreshExp,
		},
	}
}

func (s *service) Create(ctx context.Context, payload *RegisterUserPayload) (*user.User, error) {
	hash, err := hashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	//Create User Object
	user := &user.User{
		ID:         uuid.New(),
		Username:   payload.Username,
		Email:      payload.Email,
		Password:   hash,
		IsVerified: true,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, errors.New(errors.ErrInternalServer, errors.PasswordHashingFailed, err)
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.IsVerified {
		return nil, errors.New(errors.ErrInvalidRequestBody, "email not verified", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, errors.New(errors.ErrInvalidRequestBody, "wrong password", err)
	}

	return s.generateTokenPair(ctx, user.ID)
}

// Logout implements Service.
func (s *service) Logout(ctx context.Context, refreshToken string) error {
	userID, ok := m.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New(errors.ErrInternalServer, "cannot get userID from session", nil)
	}

	// 1. Get hash token via userId
	x, err := s.repo.GetRefreshTokenHash(ctx, userID)
	if err != nil {
		// Don't reveal whether token exists or not
		return errors.New(errors.ErrInternalServer, "cannot get refresh token by hash", err)
	}

	// 2. Revoke the token
	if err := s.repo.RevokeRefreshToken(ctx, x.TokenHash); err != nil {
		return errors.New(errors.ErrInternalServer, "failed to revoke token", err)
	}

	return nil
}

// RefreshTokens implements Service.
// TODO: Fix Refresh Token method
func (s *service) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// // 1. Verify the refresh token exists and isn't revoked
	// x, _ := hashPassword(refreshToken)
	// tokenHash := string(x)

	// storedToken, err := s.repo.GetRefreshTokenHash(ctx, tokenHash)
	// if err != nil {
	// 	return nil, errors.New(errors.ErrUnauthorized, "invalid refresh token", err)
	// }

	// // 2. Check if token is expired
	// if time.Now().After(storedToken.ExpiresAt) {
	// 	return nil, errors.New(errors.ErrUnauthorized, "refresh token expired", err)
	// }

	// // 3. Check if token is revoked
	// if storedToken.Revoked {
	// 	// Security measure: revoke all tokens for this user if a revoked token is used
	// 	_ = s.repo.RevokeAllRefreshTokens(ctx, storedToken.UserID)
	// 	return nil, errors.New(errors.ErrUnauthorized, "refresh token revoked", err)
	// }

	// // 4. Invalidate the current refresh token (rotation)
	// if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
	// 	return nil, errors.New(errors.ErrInternalServer, "failed to revoke token", err)
	// }

	// // 5. Generate new token pair
	// newTokens, err := s.generateTokenPair(ctx, storedToken.UserID)
	// if err != nil {
	// 	return nil, errors.New(errors.ErrInternalServer, "failed to generate tokens", err)
	// }

	return &TokenPair{}, nil
}
