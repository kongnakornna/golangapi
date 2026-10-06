package websocket

import (
	"context"
	"errors"
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
	ID          string
	UserID      string
	Role        string // RBAC role for broadcast targeting
	Conn        *websocket.Conn
	Send        chan []byte
	Hub         *Hub
	Topics      map[string]bool
	Rooms       map[string]bool
	CurrentRoom string // Currently active room (for presence tracking)
	AckManager  *AckManager

	// Sequence number tracking for message ordering
	seqMu      sync.Mutex
	nextSeq    uint64               // Next sequence number to assign
	pendingSeq map[uint64]*Envelope // Pending messages awaiting acknowledgment
	ackedSeq   uint64               // Highest acknowledged sequence number

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

// SendWithAck sends an envelope and waits for acknowledgment with retry logic.
func (c *Client) SendWithAck(ctx context.Context, env *Envelope) (*Envelope, error) {
	if c.AckManager == nil {
		return nil, errors.New("ack manager not initialized")
	}
	return c.AckManager.SendWithAck(ctx, c, env)
}

// NextSeq returns the next sequence number and increments the counter.
func (c *Client) NextSeq() uint64 {
	c.seqMu.Lock()
	defer c.seqMu.Unlock()
	c.nextSeq++
	return c.nextSeq
}

// TrackPending tracks a message awaiting acknowledgment.
func (c *Client) TrackPending(env *Envelope) {
	c.seqMu.Lock()
	defer c.seqMu.Unlock()
	if c.pendingSeq == nil {
		c.pendingSeq = make(map[uint64]*Envelope)
	}
	c.pendingSeq[env.Seq] = env
}

// HandleSeqAck processes an incoming sequence acknowledgment.
func (c *Client) HandleSeqAck(seqAck uint64) {
	c.seqMu.Lock()
	defer c.seqMu.Unlock()
	if seqAck > c.ackedSeq {
		c.ackedSeq = seqAck
	}
	// Clean up acknowledged messages
	for seq := range c.pendingSeq {
		if seq <= seqAck {
			delete(c.pendingSeq, seq)
		}
	}
}

// GetPendingCount returns the number of unacknowledged messages.
func (c *Client) GetPendingCount() int {
	c.seqMu.Lock()
	defer c.seqMu.Unlock()
	return len(c.pendingSeq)
}

// GetAckedSeq returns the highest acknowledged sequence number.
func (c *Client) GetAckedSeq() uint64 {
	c.seqMu.Lock()
	defer c.seqMu.Unlock()
	return c.ackedSeq
}

// SetCurrentRoom sets the client's current room and updates hub presence tracking.
func (c *Client) SetCurrentRoom(room string) {
	oldRoom := c.CurrentRoom
	c.CurrentRoom = room
	if c.Hub != nil && oldRoom != room {
		c.Hub.UpdateClientRoom(c, oldRoom, room)
	}
}

// GetCurrentRoom returns the client's current room.
func (c *Client) GetCurrentRoom() string {
	return c.CurrentRoom
}
