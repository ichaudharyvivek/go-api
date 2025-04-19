package v1

import (
	"context"
	"fmt"
	"net/http"

	"example.com/goapi/internal/common/query"
	"example.com/goapi/internal/domain/feed"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type FeedHandler struct {
	service   feed.Service
	validator *validator.Validate
}

func NewFeedHandler(s feed.Service, v *validator.Validate) *FeedHandler {
	return &FeedHandler{service: s, validator: v}
}

// RegisterRoutes
func (h *FeedHandler) RegisterFeedRoutes(r chi.Router) {
	r.Route("/feed", func(r chi.Router) {
		r.Get("/", h.ListUserFeed)
	})
}

func (h *FeedHandler) ListUserFeed(w http.ResponseWriter, r *http.Request) {
	// Pagination and query
	fq := &query.QueryParams{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		httpx.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fake user ID for now
	userIDStr := "22f88b40-3513-4ebe-b4ae-6363dd939455"

	// Add user ID to context
	ctx := context.WithValue(r.Context(), "user_id", userIDStr)

	// Call service
	posts, err := h.service.List(ctx, fq)
	if err != nil {
		httpx.Error(w, fmt.Sprintf("Failed to fetch feed: %s", err), http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, posts.ToDto())
}
