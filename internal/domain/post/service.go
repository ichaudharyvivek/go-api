package post

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// This is what the business layer of Posts is capable off
type Service interface {
	Create(ctx context.Context, input Form) (*Post, error)
	GetAll(ctx context.Context) (Posts, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, input Form) (*Post, error) {
	if input.Title == "" || input.Content == "" {
		return nil, errors.New("title or content cannot be empty")
	}

	p := input.ToModel()
	p.ID = uuid.New()

	err := s.repo.Save(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) GetAll(ctx context.Context) (Posts, error) {
	return s.repo.FindAll(ctx)
}
