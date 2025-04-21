package auth

import (
	"context"
	"time"

	"example.com/goapi/internal/domain/user"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *user.User) error
	GetByEmail(ctx context.Context, email string) (*user.User, error)

	// Token methods
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
