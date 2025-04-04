package router

import (
	"github.com/go-chi/chi/v5"

	"example.com/api-go/cmd/api/resource/book"
	"example.com/api-go/cmd/api/resource/health"
	_ "example.com/api-go/docs"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func New(db *gorm.DB, v *validator.Validate) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", health.Read)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/v1", func(r chi.Router) {
		bookAPI := book.New(db, v)
		r.Get("/books", bookAPI.List)
		r.Post("/books", bookAPI.Create)
		r.Get("/books/{id}", bookAPI.Get)
		r.Put("/books/{id}", bookAPI.Update)
		r.Delete("/books/{id}", bookAPI.Delete)
	})

	return r
}
