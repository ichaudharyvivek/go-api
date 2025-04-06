package post

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (Posts, error) {
	posts := make([]*Post, 0)
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repository) Create(post *Post) (*Post, error) {
	if err := r.db.Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *Repository) Read(id uuid.UUID) (*Post, error) {
	post := &Post{}
	if err := r.db.Where("id=?", id).First(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&Post{})
	return result.RowsAffected, result.Error
}

func (r *Repository) Update(post *Post) (int64, error) {
	result := r.db.Model(&Post{}).
		Select("Title", "Content", "Author", "CreatedAt", "UpdatedAt").
		Where("id = ?", post.ID).
		Updates(post)

	return result.RowsAffected, result.Error
}
