package router

import (
	"github.com/go-chi/chi/v5"

	"example.com/api-go/cmd/api/resource/book"
	"example.com/api-go/cmd/api/resource/health"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func New(db *gorm.DB, v *validator.Validate) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", health.Read)

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
