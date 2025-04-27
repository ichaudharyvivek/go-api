package auth

import (
	"context"
	"time"

	"example.com/goapi/internal/common/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Private helper methods
func (s *service) generateTokenPair(ctx context.Context, userID uuid.UUID) (*TokenPair, error) {
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Private helper methods
func (s *service) generateAccessToken(userID uuid.UUID) (string, error) {
	claims := JWTClaim{
		UserID: userID.String(),
		Roles:  []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessToken.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.accessToken.Secret))
	if err != nil {
		return "", errors.New(errors.ErrTokenGeneration, "cannot generate token", err)
	}

	return signedToken, nil
}

// Private helper methods
func (s *service) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	// Generate a secure random token
	token := uuid.New().String()

	// Hash the token before storing it
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New(errors.ErrTokenGeneration, "cannot generate token", err)
	}

	// Store in database
	if err := s.repo.CreateRefreshToken(
		ctx,
		userID,
		string(hashedToken),
		time.Now().Add(s.refreshToken.Expiration),
	); err != nil {
		return "", errors.New(errors.ErrTokenGeneration, "cannot generate token", err)
	}

	return token, nil
}

// Store password related functionality in a separate file [IMP]
func hashPassword(text string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}
