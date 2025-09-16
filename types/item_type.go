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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemTypesListResponse struct {
	ItemTypes  []ItemTypeResponse `json:"item_types"`
	TotalCount int64              `json:"total_count"`
	Limit      int                `json:"limit"`
	Offset     int                `json:"offset"`
}

type ItemTypeSearchResponse struct {
	ItemTypes   []ItemTypeResponse `json:"item_types"`
	SearchTerm  string             `json:"search_term"`
	ResultCount int                `json:"result_count"`
	Limit       int                `json:"limit"`
	Offset      int                `json:"offset"`
}
