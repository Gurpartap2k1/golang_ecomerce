package products

import (
	"gary/ecom/internal/json"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.service.ListProducts(r.Context())

	if err != nil {
		h.logger.Error("Failed to fetch products", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, 200, products)
}

func (h *Handler) ListProductById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn(
			"invalid product id",
			slog.String("id", idStr),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.service.ListProductById(r.Context(), id)
	if err != nil {
		h.logger.Error(
			"failed to fetch product",
			slog.Int64("product_id", id),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}
