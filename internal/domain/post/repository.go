package post

import "context"

// Repository defines the data access methods for Posts.
type Repository interface {
	Save(ctx context.Context, p *Post) error
	FindAll(ctx context.Context) (Posts, error)
}
