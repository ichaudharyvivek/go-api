package post

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostHandler struct {
	repository *Repository
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{
		repository: NewRepository(db),
	}
}

func (h *PostHandler) RegisterRoutes(r chi.Router) {
	r.Route("/post", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Read)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.repository.List()
	if err != nil {
		// Handle error later
		return
	}

	if len(posts) == 0 {
		if err := json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data":    []any{},
		}); err != nil {
			return
		}

		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "data": posts.ToDto()}); err != nil {
		return
	}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		return
	}

	newPost := form.ToModel()
	newPost.ID = uuid.New()

	_, err := h.repository.Create(newPost)
	if err != nil {
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "data": newPost}); err != nil {
		return
	}
}

func (h *PostHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return
	}

	post, err := h.repository.Read(id)
	if err != nil {
		return
	}

	dto := post.ToDto()
	if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "data": dto}); err != nil {
		return
	}
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		return
	}

	post := form.ToModel()
	post.ID = id

	rows, err := h.repository.Update(post)
	if err != nil || rows == 0 {
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "msg": "Resource updated successfully."}); err != nil {
		return
	}
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return
	}

	rows, err := h.repository.Delete(id)
	if err != nil || rows == 0 {
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "msg": "Resource deleted successfully."}); err != nil {
		return
	}
}
