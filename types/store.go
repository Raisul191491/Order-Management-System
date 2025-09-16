package types

import "time"

type StoreCreateRequest struct {
	Name         string `json:"name" binding:"required" validate:"required,min=1,max=255"`
	ContactPhone string `json:"contact_phone" binding:"required" validate:"required,regexp=^(01)[3-9]{1}[0-9]{8}$"`
	Address      string `json:"address"`
}

type StoreUpdateRequest struct {
	ID           int64   `json:"id" binding:"required" validate:"required,min=1"`
	Name         string  `json:"name" validate:"required,min=1,max=255"`
	ContactPhone *string `json:"contact_phone" validate:"omitempty,regexp=^(01)[3-9]{1}[0-9]{8}$"`
	Address      string  `json:"address"`
}

type StoreResponse struct {
	ID           int64     `json:"id,omitempty"`
	Name         string    `json:"name"`
	ContactPhone string    `json:"contact_phone"`
	Address      string    `json:"address,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}
