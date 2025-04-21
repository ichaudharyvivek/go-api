package repository

import (
	"context"
	"errors"
	"time"

	"example.com/goapi/internal/domain/auth"
	"example.com/goapi/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) auth.Repository {
	return &AuthRepository{db: db}
}

// User methods
func (r *AuthRepository) Create(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Token methods
func (r *AuthRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	token := &auth.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	return r.db.WithContext(ctx).Create(token).Error
}

func (r *AuthRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*auth.RefreshToken, error) {
	var token *auth.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(token).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return token, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).
		Model(&auth.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).
		Error
}

func (r *AuthRepository) RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&auth.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).
		Error
}
