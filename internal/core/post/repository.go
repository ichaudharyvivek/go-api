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

func (r *Repository) List() (Posts, error) { // Return type should be a slice of Post
	var posts Posts
	// Logic to fetch posts can be added here
	return posts, nil
}

func (r *Repository) Create(post *Post) (*Post, error) {
	// Logic to create post can be added here
	return post, nil
}

func (r *Repository) Get(id uuid.UUID) (*Post, error) {
	var post Post
	// Logic to get a post by ID can be added here
	return &post, nil
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	// Logic to delete a post by ID can be added here
	return 0, nil
}

func (r *Repository) Update(post *Post) (int64, error) {
	// Logic to update a post can be added here
	return 0, nil
}
