package feed

import (
	"context"

	"example.com/goapi/internal/common/query"
	"example.com/goapi/internal/domain/post"
	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context, userID uuid.UUID, fq *query.QueryParams) (post.Posts, error)
}
