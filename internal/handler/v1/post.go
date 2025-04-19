//	@BasePath	/api/v1

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

// ListAllPosts godoc
//
//	@Summary		List all posts
//	@Description	Get list of posts with optional filters
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			user_id	query		string		false	"User ID"
//	@Param			q		query		string		false	"Search query"
//	@Param			tags	query		[]string	false	"Tags filter (a,b,c means OR)"
//	@Param			title	query		string		false	"Exact title match"
//	@Success		200		{array}		post.DTO
//	@Failure		500		{object}	httpx.APIResponse
//	@Router			/posts [get]
func (h *Handler) ListAllPosts(w http.ResponseWriter, r *http.Request) {
	uID := uuid.Nil
	if idx := r.URL.Query().Get("user_id"); idx != "" {
		uID = uuid.MustParse(idx)
	}

	query := &post.SearchQuery{
		Query:  r.URL.Query().Get("q"),
		Tags:   r.URL.Query()["tags"],
		Title:  r.URL.Query().Get("title"),
		UserID: uID,
	}

	posts, err := h.service.FindAll(r.Context(), query)
	if err != nil {
		httpx.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, posts.ToDto())
}

// CreatePost godoc
//
//	@Summary		Create a post
//	@Description	Create a new post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		post.Form	true	"Post body"
//	@Success		201		{object}	post.DTO
//	@Failure		400		{object}	httpx.APIResponse
//	@Failure		422		{object}	httpx.APIResponse
//	@Failure		500		{object}	httpx.APIResponse
//	@Router			/posts [post]
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

	httpx.Created(w, post.ToDto())
}

// GetPostById godoc
//
//	@Summary		Get post by ID
//	@Description	Get a single post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	post.DTO
//	@Failure		400	{object}	httpx.APIResponse
//	@Failure		500	{object}	httpx.APIResponse
//	@Router			/posts/{id} [get]
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

// UpdatePostById godoc
//
//	@Summary		Update post
//	@Description	Update a post by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Post ID"
//	@Param			post	body		post.Form	true	"Updated post body"
//	@Success		200		{object}	post.DTO
//	@Failure		400		{object}	httpx.APIResponse
//	@Failure		422		{object}	httpx.APIResponse
//	@Failure		500		{object}	httpx.APIResponse
//	@Router			/posts/{id} [put]
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
		httpx.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, created.ToDto())
}

// DeletePostBy godoc
//
//	@Summary		Delete post
//	@Description	Delete a post by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{string}	string	"Deleted message"
//	@Failure		400	{object}	httpx.APIResponse
//	@Failure		500	{object}	httpx.APIResponse
//	@Router			/posts/{id} [delete]
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
