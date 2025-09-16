package types

import "time"

type DeliveryTypeCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type DeliveryTypeUpdateRequest struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type DeliveryTypeResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
