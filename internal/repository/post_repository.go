package repository

import (
	"context"

	"example.com/goapi/internal/domain/post"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) post.Repository {
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

func (r *PostRepository) FindById(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	post := &post.Post{}
	if err := r.db.WithContext(ctx).Where("id=?", id).First(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) UpdateById(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *PostRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	return nil
}
