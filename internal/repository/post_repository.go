package repository

import (
	"context"

	"example.com/goapi/internal/common/errors"
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
	err := r.db.Preload("User").WithContext(ctx).Find(&posts).Error
	return posts, err
}

func (r *PostRepository) FindById(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	post := &post.Post{}
	if err := r.db.Preload("User").WithContext(ctx).Where("id=?", id).First(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, input *post.Post) (*post.Post, error) {
	// Fetch the current post by ID to get all current data including version
	var toUpdate post.Post
	if err := r.db.WithContext(ctx).First(&toUpdate, "id = ?", input.ID).Error; err != nil {
		return nil, err
	}

	// Store the current version for optimistic locking
	currentVersion := toUpdate.Version

	// Update the fields from input (only update fields that should be updated)
	toUpdate.Title = input.Title
	toUpdate.Content = input.Content
	toUpdate.Tags = input.Tags
	// Update any other fields that need updating...

	// Increment the version for concurrency control
	toUpdate.Version = currentVersion + 1

	// Perform the update with version check
	result := r.db.WithContext(ctx).Model(&post.Post{}).
		Where("id = ? AND version = ?", toUpdate.ID, currentVersion).
		Updates(&toUpdate)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, errors.New(errors.ErrDBAccessFailure, errors.DBDataUpdateFailure, result.Error)
	}

	// Return the updated post
	return &toUpdate, nil
}

func (r *PostRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&post.Post{})
	if result.RowsAffected > 0 {
		return nil
	}

	return result.Error
}
