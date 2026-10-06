// Package http จัดเตรียม HTTP layer สำหรับโมดูล kafka (orders)
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
//
// @Summary      สร้างคำสั่งซื้อใหม่ (async ผ่าน Kafka)
// @Description  ยื่นคำสั่งซื้อ — ระบบจะ publish ไปยัง Kafka และประมวลผลแบบ async
// @Description  ผลลัพธ์ทันทีคือ 202 Accepted พร้อม order_id สำหรับติดตามสถานะ
// @Description  Client สามารถ subscribe WebSocket (`/ws?room=<order_id>`) เพื่อรับ update
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request body models.CreateOrderRequest true "Order details"
// @Success      202 {object} models.CreateOrderAcceptedResponse "Order accepted for processing"
// @Failure      400 {object} models.ErrorResponse "Invalid request body or missing required fields"
// @Failure      500 {object} models.ErrorResponse "Internal server error (Kafka publish failed)"
// @Security     BearerAuth
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ProductID == "" || req.Quantity <= 0 {
		writeJSONError(w, http.StatusBadRequest, "product_id and quantity (>=1) are required")
		return
	}

	order, err := h.usecase.CreateOrder(r.Context(), &req)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create order: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(models.CreateOrderAcceptedResponse{
		OrderID: order.ID.String(),
		Status:  order.Status,
		Message: "Order accepted for processing",
	})
}

// ============================================================
// Helpers
// ============================================================

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: msg})
}
