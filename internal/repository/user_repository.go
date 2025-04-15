package repository

import (
	"context"
	"fmt"

	"example.com/goapi/internal/common/errors"
	"example.com/goapi/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) List(ctx context.Context) (user.Users, error) {
	var users user.Users
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	user := &user.User{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrUserNotFound, fmt.Sprintf(errors.UserNotFound, id), err)
		}
		return nil, errors.New(errors.ErrInternalServer, "something went wrong", err)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *user.User) (*user.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrUserNotFound, fmt.Sprintf(errors.UserNotFound, user.ID), err)
		}
		return nil, errors.New(errors.ErrInternalServer, "something went wrong", err)
	}

	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&user.User{}).Error; err != nil {
		return err
	}
	return nil
}
