package auth

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /auth/register", h.RegisterUser)
	mux.HandleFunc("POST /auth/login", h.Login)

}
