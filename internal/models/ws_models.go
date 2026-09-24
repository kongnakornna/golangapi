package models

import (
	"net"
	"time"
)

type WsMessage struct {
	ID        uint      `gorm:"primaryKey"`
	Room      string    `gorm:"type:varchar(100);not null;index:idx_ws_messages_room"`
	Event     string    `gorm:"type:varchar(50);not null"`
	Data      string    `gorm:"type:jsonb"` // หรือใช้ datatypes.JSON
	Sender    string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_ws_messages_created_at"`
}

type WsSession struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         string    `gorm:"type:varchar(100);not null"`
	Room           string    `gorm:"type:varchar(100)"`
	ConnectedAt    time.Time `gorm:"autoCreateTime"`
	DisconnectedAt *time.Time
	ClientIP       net.IP `gorm:"type:inet"`
}
