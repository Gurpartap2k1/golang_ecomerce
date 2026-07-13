package orders

import (
	"gary/ecom/internal/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {
		h.logger.Warn(
			"Invalid Request",
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)

	if err != nil {
		h.logger.Error(
			"failed to create order",
			slog.Any("error", err),
		)
		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Order Placed", "OrderId", createdOrder.ID)

	json.Write(w, http.StatusCreated, createdOrder)
}
