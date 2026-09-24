package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// PongWait returns the pong timeout used for read deadlines.
func PongWait() time.Duration { return pongWait }

type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	Topics map[string]bool
	Rooms  map[string]bool

	sendMu sync.Mutex
	closed bool
}

func NewClient(id, userID string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:     id,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
		Topics: make(map[string]bool),
		Rooms:  make(map[string]bool),
	}
}

// safeSend enqueues a message without blocking. It returns false when the
// client has already been closed (and the Send channel closed), so callers
// never panic with "send on closed channel".
func (c *Client) safeSend(data []byte) bool {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	if c.closed {
		wsMessagesDroppedTotal.Inc()
		return false
	}
	select {
	case c.Send <- data:
		return true
	default:
		wsMessagesDroppedTotal.Inc()
		return false
	}
}

// close marks the client as closed and closes the Send channel once.
func (c *Client) close() {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	close(c.Send)
}

// WritePump sends messages from Send channel to websocket, plus periodic pings.
func (c *Client) WritePump() {
	defer c.Conn.Close()
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			wsMessagesSentTotal.Inc()
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage is a helper to send a message directly
func (c *Client) SendMessage(data []byte) {
	c.safeSend(data)
}
