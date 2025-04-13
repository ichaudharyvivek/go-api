package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/goapi/internal/common/errors"
	"example.com/goapi/internal/domain/user"
	_v "example.com/goapi/internal/utils/validator"
	"example.com/goapi/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	users, err := h.service.List(r.Context())
	if err != nil {
		fmt.Printf("Trace: %+v\n", err)
		httpx.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, users.ToDto())
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload = &user.Form{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		httpx.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		respBody := _v.ToErrResponse(err)
		httpx.Errors(w, respBody, http.StatusUnprocessableEntity)
		return
	}

	user, err := h.service.Create(r.Context(), payload)
	if err != nil {
		if apiErr, ok := err.(*errors.ApiError); ok {
			fmt.Printf("Trace: %+v\n", err)
			httpx.Error(w, apiErr.Message, http.StatusBadRequest)
			return
		}

		httpx.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.Created(w, user.ToDto())
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if apiErr, ok := err.(*errors.ApiError); ok {
			fmt.Printf("Trace: %+v\n", err)
			httpx.Error(w, apiErr.Message, http.StatusBadRequest)
			return
		}

		httpx.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, user.ToDto())
}

func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w, "UpdateUserById handler")
}

func (h *UserHandler) DeleteuserById(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w, "DeleteuserById handler")
}
