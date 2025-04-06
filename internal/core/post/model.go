package post

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DTO represents the data transfer object for a Post
type DTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Form represents the input structure for creating or updating a Post
type Form struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required,alphaspace,max=255"`
}

// Post represents the database model for a Post
type Post struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Title     string
	Content   string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// Posts represents a collection of posts
type Posts []*Post

// ToModel converts a Form into a Post model
func (f *Form) ToModel() *Post {
	return &Post{
		Title:   f.Title,
		Content: f.Content,
		Author:  f.Author,
	}
}

// ToDto converts a Post model into a DTO
func (p *Post) ToDto() *DTO {
	return &DTO{
		ID:        p.ID.String(),
		Title:     p.Title,
		Content:   p.Content,
		Author:    p.Author,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToDto converts a collection of Post models into a slice of DTOs
func (items Posts) ToDto() []*DTO {
	dtos := make([]*DTO, len(items))
	for i, v := range items {
		dtos[i] = v.ToDto()
	}
	return dtos
}
