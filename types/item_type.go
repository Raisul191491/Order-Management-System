package types

import "time"

type ItemTypeCreateRequest struct {
	Name string `json:"name" binding:"required" validate:"required,max=50"`
}

type ItemTypeUpdateRequest struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,max=50"`
}

type ItemTypeResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ItemTypesListResponse struct {
	ItemTypes  []ItemTypeResponse `json:"itemTypes"`
	TotalCount int64              `json:"totalCount"`
	Limit      int                `json:"limit"`
	Offset     int                `json:"offset"`
}
