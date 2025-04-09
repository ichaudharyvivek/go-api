package v1

import (
	"encoding/json"
	"net/http"

	"example.com/goapi/internal/domain/post"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service post.Service
}

func NewHandler(s post.Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes mounts the post routes on the given router
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
	})
}

// Create handles POST /posts
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input post.Form
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.service.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetAll handles GET /posts
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
