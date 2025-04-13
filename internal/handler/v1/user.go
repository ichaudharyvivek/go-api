package v1

import (
	"net/http"

	"example.com/goapi/internal/domain/user"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service   user.Service
	validator *validator.Validate
}

func NewUserHandler(s user.Service, v *validator.Validate) *UserHandler {
	return &UserHandler{service: s, validator: v}
}

// RegisterRoutes mounts the post routes on the given router
func (h *UserHandler) RegisterUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.ListAllUsers)
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUserById)
		r.Put("/{id}", h.UpdateUserById)
		r.Delete("/{id}", h.DeleteuserById)
	})
}

func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w, "ListAllUsers handler")
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w, "CreateUser handler")
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w, "GetUserById handler")
}

func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w, "UpdateUserById handler")
}

func (h *UserHandler) DeleteuserById(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w, "DeleteuserById handler")
}
