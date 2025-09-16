package types

import "time"

type ZoneCreateRequest struct {
	CityID int64  `json:"cityId" binding:"required" validate:"required"`
	Name   string `json:"name" binding:"required" validate:"required,max=100"`
}

// ZoneUpdateRequest represents the request structure for updating an existing zone
type ZoneUpdateRequest struct {
	ID   int64  `json:"id" binding:"required" validate:"required"`
	Name string `json:"name" validate:"required,max=100"`
}

// ZoneResponse represents the response structure for zone data
type ZoneResponse struct {
	ID        int64     `json:"id"`
	CityID    int64     `json:"cityId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ZonesListResponse represents the response structure for paginated zone lists
type ZonesListResponse struct {
	Zones      []ZoneResponse `json:"zones"`
	TotalCount int64          `json:"total_count"`
	Limit      int            `json:"limit"`
	Offset     int            `json:"offset"`
}
