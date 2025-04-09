package repository

import (
	"context"

	"example.com/goapi/internal/domain/post"
	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, p *post.Post) error
	FindAll(ctx context.Context) (post.Posts, error)
}

type PostRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, p *post.Post) error {
	return r.db.WithContext(ctx).Create(p).Error

}

func (r *PostRepository) FindAll(ctx context.Context) (post.Posts, error) {
	var posts post.Posts
	err := r.db.WithContext(ctx).Find(&posts).Error
	return posts, err
}
