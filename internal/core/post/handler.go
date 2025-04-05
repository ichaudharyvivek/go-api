package post

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	r.Route("/v1", func(r chi.Router) {
		r.Get("/books", h.List)
		r.Post("/books", h.Create)
		r.Get("/books/{id}", h.Read)
		r.Put("/books/{id}", h.Update)
		r.Delete("/books/{id}", h.Delete)
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

	if err := json.NewEncoder(w).Encode(posts.ToDto()); err != nil {
		return
	}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new book"))
}

func (h *PostHandler) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Get book with ID: " + id))
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Update book with ID: " + id))
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Delete book with ID: " + id))
}
