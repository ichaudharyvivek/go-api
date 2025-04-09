package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	e "example.com/goapi/internal/common/err"
	"example.com/goapi/internal/domain/post"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service post.Service
}

func NewHandler(s post.Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes mounts the post routes on the given router
func (h *Handler) RegisterRoutes(r chi.Router, v *validator.Validate) {
	r.Route("/post", func(r chi.Router) {
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
	})
}

// Create handles POST /posts
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input post.Form
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, http.StatusBadRequest, e.RespJSONDecodeFailure)
		return
	}

	post, err := h.service.Create(r.Context(), input)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, e.RespDBDataInsertFailure)
		return
	}

	httpx.Created(w, post)
}

// GetAll handles GET /posts
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAll(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		return
	}

	httpx.Ok(w, posts.ToDto())
}
