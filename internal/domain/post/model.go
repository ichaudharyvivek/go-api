package post

import (
	"time"

	"example.com/goapi/internal/domain/user"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Post represents the database model for a Post
type Post struct {
	ID        uuid.UUID      `gorm:"primarykey"`
	Title     string         `gorm:"type:varchar(255);not null"`
	Content   string         `gorm:"type:text;not null"`
	Tags      pq.StringArray `gorm:"type:text[]"`
	UserID    *uuid.UUID     `gorm:"type:uuid;index"`
	User      *user.User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Posts represents a collection of posts
type Posts []*Post

// DTO represents the data transfer object for a Post
type DTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	User      *user.DTO `json:"user,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// Form represents the input structure for creating or updating a Post
type Form struct {
	Title   string   `json:"title" validate:"required,max=255"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
	UserID  string   `json:"user_id" validate:"required"`
}

// ToModel converts a Form into a Post model
func (f *Form) ToModel() *Post {
	userID := uuid.MustParse(f.UserID)

	return &Post{
		ID:      uuid.New(),
		Title:   f.Title,
		Content: f.Content,
		Tags:    f.Tags,
		UserID:  &userID,
	}
}

// ToDto converts a Post model into a DTO
func (p *Post) ToDto() *DTO {
	var userDto *user.DTO
	if p.UserID != nil && p.User != nil && p.User.ID != uuid.Nil {
		userDto = p.User.ToDto()
	}

	return &DTO{
		ID:        p.ID.String(),
		Title:     p.Title,
		Content:   p.Content,
		Tags:      p.Tags,
		User:      userDto,
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
