package router

import (
	"expvar"
	"net/http"
	"time"

	"example.com/goapi/internal/database/cache"
	"example.com/goapi/internal/domain/auth"
	"example.com/goapi/internal/domain/feed"
	"example.com/goapi/internal/domain/post"
	"example.com/goapi/internal/domain/user"
	v1 "example.com/goapi/internal/handler/v1"
	m "example.com/goapi/internal/middleware"
	"example.com/goapi/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"

	_ "example.com/goapi/docs"
)

func NewRouter(db *gorm.DB, v *validator.Validate, rd *cache.Client) http.Handler {
	r := chi.NewRouter()
	applyMiddlewares(r)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/debug/vars", expvar.Handler().ServeHTTP)
	r.Route("/api/v1", func(r chi.Router) {
		registerPostRoutes(r, db, v, rd)
		registerUserRoutes(r, db, v, rd)
		registerFeedRoutes(r, db, v, rd)
		registerAuthRoutes(r, db, v, rd)
	})

	return r
}

func registerPostRoutes(r chi.Router, db *gorm.DB, v *validator.Validate, rd *cache.Client) {
	repo := repository.NewRepository(db)
	service := post.NewService(repo)
	handler := v1.NewHandler(service, v, rd)
	handler.RegisterRoutes(r)
}

func registerUserRoutes(r chi.Router, db *gorm.DB, v *validator.Validate, rd *cache.Client) {
	repo := repository.NewUserRepository(db)
	service := user.NewService(repo)
	handler := v1.NewUserHandler(service, v)
	handler.RegisterUserRoutes(r)
}

func registerFeedRoutes(r chi.Router, db *gorm.DB, v *validator.Validate, rd *cache.Client) {
	repo := repository.NewFeedRepository(db)
	service := feed.NewService(repo)
	handler := v1.NewFeedHandler(service, v)
	handler.RegisterFeedRoutes(r)
}

func registerAuthRoutes(r chi.Router, db *gorm.DB, v *validator.Validate, rd *cache.Client) {
	repo := repository.NewAuthRepository(db)
	service := auth.NewService(repo, "secret", "refresh", 15*time.Minute, 7*24*time.Hour)
	handler := v1.NewAuthHandler(service, v)
	handler.RegisterAuthRoutes(r)
}

func applyMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	// r.Use(middleware.Logger) // Not required because we have our own middleware for logging
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(m.LogContext)
	r.Use(cors.Handler(cors.Options{
		// Allow all origins for simplicity
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		ExposedHeaders: []string{"X-Custom-Header"},
		MaxAge:         300, // 5 minutes
	}))
}
