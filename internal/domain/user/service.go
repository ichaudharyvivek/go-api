package user

import (
	"context"

	"github.com/google/uuid"
)

// This is what the business layer of Users is capable off
type Service interface {
	List(ctx context.Context) (Users, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) List(ctx context.Context) (Users, error) {
	return nil, nil
}

func (s *service) Create(ctx context.Context, user *User) error {
	return nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return nil, nil
}

func (s *service) Update(ctx context.Context, user *User) (*User, error) {
	return nil, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
