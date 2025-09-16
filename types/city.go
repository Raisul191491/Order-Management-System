package types

import "time"

type CityCreateRequest struct {
	Name            string  `json:"name" binding:"required"`
	BaseDeliveryFee float64 `json:"baseDeliveryFee"`
}

type CityUpdateRequest struct {
	ID              int64   `json:"id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	BaseDeliveryFee float64 `json:"baseDeliveryFee"`
}

type CityResponse struct {
	Id              int64     `json:"id,omitempty"`
	Name            string    `json:"name"`
	BaseDeliveryFee float64   `json:"baseDeliveryFee,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
}
