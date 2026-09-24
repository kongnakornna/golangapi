package websocket

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	wsConnectedClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ws_connected_clients",
		Help: "Current number of connected WebSocket clients",
	})
	wsConnectionsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ws_connections_total",
		Help: "Total number of WebSocket connections established",
	})
	wsActiveRooms = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ws_active_rooms",
		Help: "Current number of active WebSocket rooms",
	})
	wsActiveTopics = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ws_active_topics",
		Help: "Current number of active WebSocket topic subscriptions",
	})
	wsMessagesBroadcastTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ws_messages_broadcast_total",
		Help: "Total number of messages delivered to WebSocket clients per broadcast type",
	}, []string{"type"})
	wsMessagesReceivedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ws_messages_received_total",
		Help: "Total number of frames received from WebSocket clients",
	}, []string{"type"})
	wsMessagesSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ws_messages_sent_total",
		Help: "Total number of messages written to WebSocket clients",
	})
	wsMessagesDroppedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ws_messages_dropped_total",
		Help: "Total number of messages dropped because a client send buffer was full or closed",
	})
)

// RecordReceivedMessage increments the received-frame counter for the given
// frame type (subscribe, unsubscribe, join_room, leave_room, message, unknown).
func RecordReceivedMessage(kind string) {
	if kind == "" {
		kind = "unknown"
	}
	wsMessagesReceivedTotal.WithLabelValues(kind).Inc()
}
