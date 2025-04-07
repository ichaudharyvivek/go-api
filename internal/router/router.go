package rotuer

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		// healthHandler := health.NewHealthHandler()
		// healthHandler.RegisterRoutes(r)

		// postHandler := post.NewPostHandler(db)
		// postHandler.RegisterRoutes(r)
	})

	return r
}
