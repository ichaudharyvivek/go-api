package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents the database model for a user.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"size:255;not null;unique"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  []byte    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Users represent the list of user
type Users []*User

// DTO represents the data transfer object for User
type DTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Form represent the input structure for creating user or updating user
type Form struct {
	Username string `json:"username" validate:"required,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// ToModel converts a form to user model
func (f *Form) ToModel() *User {
	return &User{
		ID:       uuid.New(),
		Username: f.Username,
		Email:    f.Email,
		Password: []byte(f.Password),
	}
}

// ToDto converts a User model to a DTO
func (u *User) ToDto() *DTO {
	return &DTO{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToDto converts a collection of User models to slice of DTOs
func (list Users) ToDto() []*DTO {
	dtos := make([]*DTO, len(list))
	for i, v := range list {
		dtos[i] = v.ToDto()
	}

	return dtos
}
