//	@BasePath	/api/v1

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
		r.Delete("/{id}", h.DeleteUserById)
	})
}

// ListAllUsers godoc
//	@Summary		List all users
//	@Description	Get a list of all registered users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		user.DTO
//	@Failure		500	{object}	httpx.APIResponse
//	@Router			/users [get]
func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		fmt.Printf("Trace: %+v\n", err)
		httpx.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.Ok(w, users.ToDto())
}

// CreateUser godoc
//	@Summary		Create a new user
//	@Description	Create a user with email and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		user.Form	true	"User form payload"
//	@Success		201		{object}	user.DTO
//	@Failure		400		{object}	httpx.APIResponse
//	@Failure		422		{object}	httpx.APIResponse
//	@Router			/users [post]
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

// GetUserById godoc
//	@Summary		Get user by ID
//	@Description	Fetch a single user by their UUID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	user.DTO
//	@Failure		400	{object}	httpx.APIResponse
//	@Failure		500	{object}	httpx.APIResponse
//	@Router			/users/{id} [get]
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

// UpdateUserById godoc
//	@Summary		Update user by ID
//	@Description	Update a user's information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"User ID"
//	@Param			user	body		user.Form	true	"User update payload"
//	@Success		200		{object}	user.DTO
//	@Failure		400		{object}	httpx.APIResponse
//	@Failure		422		{object}	httpx.APIResponse
//	@Router			/users/{id} [put]
func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	// TODO: Add Update route
	httpx.Ok(w, "UpdateUserById handler")
}

// DeleteUserById godoc
//	@Summary		Delete user by ID
//	@Description	Permanently remove a user by UUID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{string}	string	"Deleted message"
//	@Failure		400	{object}	httpx.APIResponse
//	@Failure		500	{object}	httpx.APIResponse
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	// TODO: Add Delete route
	httpx.Ok(w, "DeleteUserById handler")
}
