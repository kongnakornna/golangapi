package presenter

import "icmongolang/pkg/helpers"

type CreateTemplateRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateTemplateRequest struct {
	Name   string `json:"name,omitempty"`
	Status *int   `json:"status,omitempty"`
}

type TemplateResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Status    int               `json:"status"`
	CreatedAt helpers.LocalTime `json:"created_at"`
	UpdatedAt helpers.LocalTime `json:"updated_at"`
}
