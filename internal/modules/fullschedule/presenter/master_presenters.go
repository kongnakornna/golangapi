package presenter

import (
	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// ============================================================================
// Group
// ============================================================================

// GroupCreate is the payload for creating a group.
// คำขอสร้างกลุ่ม
type GroupCreate struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description,omitempty"`
	SortID      int     `json:"sort_id"`
	Status      string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

// GroupUpdate is the payload for updating a group.
// คำขอแก้ไขกลุ่ม
type GroupUpdate struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
	SortID      *int    `json:"sort_id,omitempty"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

// GroupResponse is the response shape of a group.
// ข้อมูลกลุ่มที่ส่งกลับ
type GroupResponse struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Status      string            `json:"status"`
	SortID      int               `json:"sort_id"`
	CreatedBy   *uuid.UUID        `json:"created_by,omitempty"`
	CreatedAt   helpers.LocalTime `json:"created_at"`
	UpdatedBy   *uuid.UUID        `json:"updated_by,omitempty"`
	UpdatedAt   helpers.LocalTime `json:"updated_at"`
	Version     int64             `json:"version"`
}

// PaginatedGroupResponse holds paginated groups.
// รายการกลุ่มแบบแบ่งหน้า
type PaginatedGroupResponse struct {
	Groups     []*GroupResponse `json:"groups"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

// ============================================================================
// Zone
// ============================================================================

// ZoneCreate is the payload for creating a zone (requires group_id).
// คำขอสร้างโซน
type ZoneCreate struct {
	GroupID     uuid.UUID `json:"group_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description *string   `json:"description,omitempty"`
	SortID      int       `json:"sort_id"`
}

// ZoneUpdate is the payload for updating a zone.
// คำขอแก้ไขโซน
type ZoneUpdate struct {
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string    `json:"description,omitempty"`
	SortID      *int       `json:"sort_id,omitempty"`
}

// ZoneResponse is the response shape of a zone.
// ข้อมูลโซนที่ส่งกลับ
type ZoneResponse struct {
	ID          uuid.UUID         `json:"id"`
	GroupID     uuid.UUID         `json:"group_id"`
	GroupName   string            `json:"group_name"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	SortID      int               `json:"sort_id"`
	CreatedBy   *uuid.UUID        `json:"created_by,omitempty"`
	CreatedAt   helpers.LocalTime `json:"created_at"`
	UpdatedBy   *uuid.UUID        `json:"updated_by,omitempty"`
	UpdatedAt   helpers.LocalTime `json:"updated_at"`
	Version     int64             `json:"version"`
}

// PaginatedZoneResponse holds paginated zones.
// รายการโซนแบบแบ่งหน้า
type PaginatedZoneResponse struct {
	Zones      []*ZoneResponse `json:"zones"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// ============================================================================
// Area
// ============================================================================

// AreaCreate is the payload for creating an area (requires zone_id).
// คำขอสร้างพื้นที่
type AreaCreate struct {
	ZoneID      uuid.UUID `json:"zone_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description *string   `json:"description,omitempty"`
	SortID      int       `json:"sort_id"`
}

// AreaUpdate is the payload for updating an area.
// คำขอแก้ไขพื้นที่
type AreaUpdate struct {
	ZoneID      *uuid.UUID `json:"zone_id,omitempty"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string    `json:"description,omitempty"`
	SortID      *int       `json:"sort_id,omitempty"`
}

// AreaResponse is the response shape of an area.
// ข้อมูลพื้นที่ที่ส่งกลับ
type AreaResponse struct {
	ID          uuid.UUID         `json:"id"`
	ZoneID      uuid.UUID         `json:"zone_id"`
	ZoneName    string            `json:"zone_name"`
	GroupName   string            `json:"group_name"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	SortID      int               `json:"sort_id"`
	CreatedBy   *uuid.UUID        `json:"created_by,omitempty"`
	CreatedAt   helpers.LocalTime `json:"created_at"`
	UpdatedBy   *uuid.UUID        `json:"updated_by,omitempty"`
	UpdatedAt   helpers.LocalTime `json:"updated_at"`
	Version     int64             `json:"version"`
}

// PaginatedAreaResponse holds paginated areas.
// รายการพื้นที่แบบแบ่งหน้า
type PaginatedAreaResponse struct {
	Areas      []*AreaResponse `json:"areas"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// ============================================================================
// Device <-> Area mapping
// ============================================================================

// DeviceMapRequest binds IoT devices to an area.
// คำขอแมปอุปกรณ์เข้ากับพื้นที่
type DeviceMapRequest struct {
	DeviceIDs []int `json:"device_ids"`
}

// MessageResponse is a generic {"message": "..."} payload.
// ข้อความตอบกลับแบบทั่วไป
type MessageResponse struct {
	Message string `json:"message"`
}

// AreaDeviceResponse describes a device mapped to an area.
// อุปกรณ์ที่แมปกับพื้นที่
type AreaDeviceResponse struct {
	DeviceID int     `json:"device_id"`
	DeviceSN *string `json:"device_sn,omitempty"`
}

// ============================================================================
// Scope preview
// ============================================================================

// ScopePreviewRequest selects a group/zone/area combination to preview.
// คำขอพรีวิวขอบเขตการทำงาน
type ScopePreviewRequest struct {
	GroupID *uuid.UUID `json:"group_id"`
	ZoneID  *uuid.UUID `json:"zone_id"`
	AreaID  *uuid.UUID `json:"area_id"`
}

// ScopePreviewResponse shows how many devices the scope covers.
// จำนวนอุปกรณ์ที่อยู่ในขอบเขต
type ScopePreviewResponse struct {
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
	ZoneID      *uuid.UUID `json:"zone_id,omitempty"`
	AreaID      *uuid.UUID `json:"area_id,omitempty"`
	DeviceCount int64      `json:"device_count"`
}
