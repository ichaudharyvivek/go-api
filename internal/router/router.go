package router

import (
	"net/http"
	"time"

	"example.com/goapi/internal/domain/post"
	v1 "example.com/goapi/internal/handler/v1"
	postRepo "example.com/goapi/internal/repository"
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

	r.Route("/api/v1", func(r chi.Router) {
		// healthHandler := health.NewHealthHandler()
		// healthHandler.RegisterRoutes(r)

		postRepo := postRepo.NewRepository(db)
		postService := post.NewService(postRepo)
		postHandler := v1.NewHandler(postService)
		postHandler.RegisterRoutes(r, v)
	})

	return r
}
