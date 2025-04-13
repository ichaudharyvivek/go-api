package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the methods for interacting with the User data store.
type Repository interface {
	List(ctx context.Context) (Users, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
