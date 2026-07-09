package auth

import (
	"errors"
	"gary/ecom/internal/json"
	"log"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var tempUser UserRequest
	if err := json.Read(r, &tempUser); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createUser, err := h.service.RegisterUser(r.Context(), tempUser)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createUser)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.Read(r, &req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//calling the service
	token, err := h.service.Login(r.Context(), req)

	if err != nil {
		log.Println(err)
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, LoginResponse{
		Token: token,
	})
}
