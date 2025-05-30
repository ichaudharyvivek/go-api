//	@BasePath	/api/v1

package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/goapi/internal/common/errors"
	"example.com/goapi/internal/database/cache"
	"example.com/goapi/internal/domain/post"
	m "example.com/goapi/internal/middleware"
	_v "example.com/goapi/internal/utils/validator"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	service   post.Service
	validator *validator.Validate
	redis     *cache.Client
}

func NewHandler(s post.Service, v *validator.Validate, rd *cache.Client) *Handler {
	return &Handler{service: s, validator: v, redis: rd}
}

// RegisterRoutes mounts the post routes on the given router
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", h.ListAllPosts)
		r.Get("/{id}", h.GetPostById)

		r.Group(func(r chi.Router) {
			r.Use(m.Authenticate())
			r.With(m.AllowAccess([]string{"user", "admin"})).Post("/", h.CreatePost)
			r.Put("/{id}", h.UpdatePostById)
			r.Delete("/{id}", h.DeletePostBy)
		})
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
	// Generate a cache key based on all query parameters
	cacheKey := fmt.Sprintf("posts:%s:%s:%v:%s",
		r.URL.Query().Get("user_id"),
		r.URL.Query().Get("q"),
		r.URL.Query()["tags"],
		r.URL.Query().Get("title"),
	)

	// Try to get from cache first
	cachedData, err := h.redis.Get(r.Context(), cacheKey).Bytes()
	if err == nil {
		var posts post.Posts
		err = json.Unmarshal(cachedData, &posts)
		if err != nil {
			httpx.Error(w, fmt.Sprintf("Error unmarshalling cache data: %s", err), http.StatusInternalServerError)
			return
		}

		httpx.Ok(w, posts)
		return
	}

	// Cache miss - proceed with normal processing
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

	// Convert to DTO and cache the result
	postsDto := posts.ToDto()
	jsonData, err := json.Marshal(postsDto)
	if err == nil {
		// Cache for 5 minutes (adjust TTL as needed)
		h.redis.Set(r.Context(), cacheKey, jsonData, 5*time.Minute)
	}

	httpx.Ok(w, postsDto)
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
