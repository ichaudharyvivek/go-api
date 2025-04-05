package http

import (
	"example.com/goapi/internal/handler/http/health"
	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	healthHandler := health.NewHealthHandler()
	healthHandler.RegisterRoutes(r)

	return r
}
