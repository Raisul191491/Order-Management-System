package types

import "time"

type CityCreateRequest struct {
	Name            string  `json:"name" binding:"required"`
	BaseDeliveryFee float64 `json:"base_delivery_fee"`
}

type CityUpdateRequest struct {
	ID              int64   `json:"id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	BaseDeliveryFee float64 `json:"base_delivery_fee"`
}

type CityResponse struct {
	Id              int64     `json:"id,omitempty"`
	Name            string    `json:"name"`
	BaseDeliveryFee float64   `json:"base_delivery_fee,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}
