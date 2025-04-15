package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/goapi/internal/common/errors"
	"example.com/goapi/internal/domain/post"
	_v "example.com/goapi/internal/utils/validator"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", h.ListAllPosts)
		r.Post("/", h.CreatePost)
		r.Get("/{id}", h.GetPostById)
		r.Put("/{id}", h.UpdatePostById)
		r.Delete("/{id}", h.DeletePostBy)
	})
}

// ListAllPosts handles GET /posts
func (h *Handler) ListAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.FindAll(r.Context())
	if err != nil {
		httpx.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, posts.ToDto())
}

// CreatePost handles POST /posts
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var input = &post.Form{}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		httpx.Error(w, errors.JSONDecodeFailure, http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(input); err != nil {
		respBody := _v.ToErrResponse(err)
		httpx.Errors(w, respBody, http.StatusUnprocessableEntity)
		return
	}

	post, err := h.service.Create(r.Context(), input)
	if err != nil {
		httpx.Error(w, errors.DBDataInsertFailure, http.StatusInternalServerError)
		return
	}

	httpx.Created(w, post)
}

func (h *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, errors.InvalidURLParamID, http.StatusBadRequest)
		return
	}

	post, err := h.service.FindById(r.Context(), id)
	if err != nil {
		httpx.Error(w, errors.DBDataAccessFailure, http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, post.ToDto())
}

func (h *Handler) UpdatePostById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, errors.InvalidURLParamID, http.StatusBadRequest)
		return
	}

	input := &post.Form{}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		httpx.Error(w, errors.JSONDecodeFailure, http.StatusBadRequest)
		return
	}

	post := input.ToModel()
	post.ID = id

	created, err := h.service.Update(r.Context(), post)
	if err != nil {
		httpx.Error(w, errors.DBDataInsertFailure, http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, created.ToDto())
}

func (h *Handler) DeletePostBy(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, errors.InvalidURLParamID, http.StatusBadRequest)
		return
	}

	err = h.service.DeleteById(r.Context(), id)
	if err != nil {
		httpx.Error(w, errors.DBDataRemoveFailure, http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, fmt.Sprintf("Deleted Post with id `%s`", id))
}
