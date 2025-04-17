package router

import (
	"net/http"
	"time"

	"example.com/goapi/internal/domain/post"
	"example.com/goapi/internal/domain/user"
	v1 "example.com/goapi/internal/handler/v1"
	m "example.com/goapi/internal/middleware"
	"example.com/goapi/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, v *validator.Validate) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(m.LogTime())

	r.Route("/api/v1", func(r chi.Router) {
		registerPostRoutes(r, db, v)
		registerUserRoutes(r, db, v)
	})

	return r
}

func registerPostRoutes(r chi.Router, db *gorm.DB, v *validator.Validate) {
	repo := repository.NewRepository(db)
	service := post.NewService(repo)
	handler := v1.NewHandler(service, v)
	handler.RegisterRoutes(r)
}

func registerUserRoutes(r chi.Router, db *gorm.DB, v *validator.Validate) {
	repo := repository.NewUserRepository(db)
	service := user.NewService(repo)
	handler := v1.NewUserHandler(service, v)
	handler.RegisterUserRoutes(r)
}
