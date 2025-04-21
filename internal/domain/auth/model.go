package auth

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validated:"required,min=3,max=72"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenPayload struct {
	TokenID uuid.UUID `json:"token_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type TokenConfig struct {
	Secret     string
	Expiration time.Duration
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"type:text;not null;index"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Revoked   bool      `gorm:"default:false"`
}
