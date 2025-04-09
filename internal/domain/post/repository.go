package post

import "context"

// Repository defines the data access methods for Posts.
type Repository interface {
	FindAll(ctx context.Context) (Posts, error)
	Create(ctx context.Context, p *Post) error
}
