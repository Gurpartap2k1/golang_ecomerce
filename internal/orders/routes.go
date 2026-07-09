package orders

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", h.PlaceOrder)
	})
}
