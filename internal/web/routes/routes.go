package routes

import (
	"github.com/Aneeshie/loom/internal/web/handlers"
	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router, h *handlers.Handler) {
	r.Get("/health", h.CheckHealth)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/query", h.HandleQuery)
	})
}
