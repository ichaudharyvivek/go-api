package repository

import (
	"context"
	"errors"

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

func (r *PostRepository) Update(ctx context.Context, input *post.Post) (*post.Post, error) {
	var ep *post.Post
	if err := r.db.WithContext(ctx).First(&ep, "id = ?", input.ID).Error; err != nil {
		return nil, err
	}

	result := r.db.WithContext(ctx).Model(&post.Post{}).Select("Title", "Content", "Author", "CreatedAt", "UpdatedAt").Where("id = ?", input.ID)
	if result.RowsAffected > 0 {
		return input, nil
	}

	return nil, errors.New("post not found or no changes made")
}

func (r *PostRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&post.Post{})
	if result.RowsAffected > 0 {
		return nil
	}

	return result.Error
}
