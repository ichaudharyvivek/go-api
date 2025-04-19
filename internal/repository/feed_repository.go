package repository

import (
	"context"

	"example.com/goapi/internal/common/errors"
	"example.com/goapi/internal/common/query"
	"example.com/goapi/internal/domain/feed"
	"example.com/goapi/internal/domain/post"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) feed.Repository {
	return &FeedRepository{db: db}
}

func (r *FeedRepository) List(ctx context.Context, userID uuid.UUID, fq *query.QueryParams) (post.Posts, error) {
	feed := post.Posts{}

	db := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID)
	db = query.Apply(db, fq)
	if err := db.Find(&feed).Error; err != nil {
		return nil, errors.New(errors.ErrDBAccessFailure, "Resource Not Found", err)
	}

	return feed, nil
}
