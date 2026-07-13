package products

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/products", h.ListProducts)
	mux.HandleFunc("/products/{id}", h.ListProductById)
}
