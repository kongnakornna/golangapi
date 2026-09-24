package wos

import "net/http"

// Handlers defines HTTP handler methods for the Web Order System.
// ตัวจัดการ HTTP สำหรับระบบสั่งซื้อออนไลน์
type Handlers interface {
	CreateOrder() func(w http.ResponseWriter, r *http.Request)
	GetOrder() func(w http.ResponseWriter, r *http.Request)
	ListOrders() func(w http.ResponseWriter, r *http.Request)
	UpdateOrderStatus() func(w http.ResponseWriter, r *http.Request)
}
