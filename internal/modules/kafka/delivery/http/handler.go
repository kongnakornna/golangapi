package http

import (
	"encoding/json"
	"net/http"

	"icmongolang/internal/modules/kafka/models"
	"icmongolang/internal/modules/kafka/usecase"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

// CreateOrder godoc
// @Summary      Create a new order (async via Kafka)
// @Description  Submit an order request. The order will be processed asynchronously by Kafka consumer.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request body models.CreateOrderRequest true "Order details"
// @Success      202  {object}  map[string]interface{}  "Order accepted"
// @Failure      400  {object}  map[string]string       "Invalid request"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProductID == "" || req.Quantity <= 0 {
		http.Error(w, "product_id and quantity (>=1) are required", http.StatusBadRequest)
		return
	}

	order, err := h.usecase.CreateOrder(r.Context(), &req)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": order.ID.String(),
		"status":   order.Status,
		"message":  "Order accepted for processing",
	})
}
