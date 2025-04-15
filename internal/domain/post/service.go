package post

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// This is what the business layer of Posts is capable off
type Service interface {
	Create(ctx context.Context, input *Form) (*Post, error)
	FindAll(ctx context.Context) (Posts, error)
	Update(ctx context.Context, input *Post) (*Post, error)
	FindById(ctx context.Context, id uuid.UUID) (*Post, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, input *Form) (*Post, error) {
	if input.Title == "" || input.Content == "" {
		return nil, errors.New("title or content cannot be empty")
	}

	p := input.ToModel()
	err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) FindAll(ctx context.Context) (Posts, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) FindById(ctx context.Context, id uuid.UUID) (*Post, error) {
	post, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) Update(ctx context.Context, input *Post) (*Post, error) {
	updated, err := s.repo.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *service) DeleteById(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteById(ctx, id)
}
