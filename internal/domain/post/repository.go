package post

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the data access methods for Posts.
type Repository interface {
	Create(ctx context.Context, p *Post) error
	FindAll(ctx context.Context) (Posts, error)
	FindById(ctx context.Context, id uuid.UUID) (*Post, error)
	Update(ctx context.Context, input *Post) (*Post, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}
