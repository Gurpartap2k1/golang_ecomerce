package orders

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler, auth func(http.Handler) http.Handler) {
	mux.Handle("POST /orders", auth(http.HandlerFunc(h.PlaceOrder)))
}
