package presenter

import (
	"github.com/google/uuid"
)

type ItemCreate struct {
	Title       string `json:"title" validate:"required" example:"item title"`
	Description string `json:"description" validate:"required" example:"item description"`
}

type ItemResponse struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	OwnerId     uuid.UUID `json:"owner_id,omitempty"`
}

type ItemUpdate struct {
	Title       string `json:"title" example:"item title"`
	Description string `json:"description" example:"item description"`
}

// PaginatedItemsResponse – all fields exported (uppercase)
type PaginatedItemsResponse struct {
	Items      []*ItemResponse `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}
