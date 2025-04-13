package user

import (
	"context"

	"example.com/goapi/internal/common/errors"
	"github.com/google/uuid"
)

// This is what the business layer of Users is capable off
type Service interface {
	List(ctx context.Context) (Users, error)
	Create(ctx context.Context, form *Form) (*User, error)
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
	return s.repo.List(ctx)
}

func (s *service) Create(ctx context.Context, form *Form) (*User, error) {
	user, err := form.ToModel()
	if err != nil {
		return nil, errors.New(errors.ErrInternalServer, errors.PasswordHashingFailed, err)
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, errors.New(errors.ErrDBInsertFailure, errors.DBDataInsertFailure, err)
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New(errors.ErrDBAccessFailure, errors.DBDataAccessFailure, err)
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, user *User) (*User, error) {
	user, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, errors.New(errors.ErrDBUpdateFailure, errors.DBDataUpdateFailure, err)
	}

	return user, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.New(errors.ErrDBDeleteFailure, errors.DBDataRemoveFailure, err)
	}

	return nil
}
