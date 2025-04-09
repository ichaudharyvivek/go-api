package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	e "example.com/goapi/internal/common/err"
	"example.com/goapi/internal/domain/post"
	_v "example.com/goapi/internal/utils/validator"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service   post.Service
	validator *validator.Validate
}

func NewHandler(s post.Service, v *validator.Validate) *Handler {
	return &Handler{service: s, validator: v}
}

// RegisterRoutes mounts the post routes on the given router
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/post", func(r chi.Router) {
		r.Get("/", h.FindAll)
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

	if err := h.validator.Struct(input); err != nil {
		respBody := _v.ToErrResponse(err)
		httpx.Errors(w, http.StatusUnprocessableEntity, respBody)
		return
	}

	post, err := h.service.Create(r.Context(), input)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, e.RespDBDataInsertFailure)
		return
	}

	httpx.Created(w, post)
}

// FindAll handles GET /posts
func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.FindAll(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		return
	}

	httpx.Ok(w, posts.ToDto())
}
