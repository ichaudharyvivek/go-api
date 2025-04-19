package feed

import (
	"context"

	"example.com/goapi/internal/common/query"
	"example.com/goapi/internal/domain/post"
	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context, fq *query.QueryParams) (post.Posts, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) List(ctx context.Context, fq *query.QueryParams) (post.Posts, error) {
	uID := uuid.MustParse(ctx.Value("user_id").(string))
	return s.repo.List(ctx, uID, fq)
}
